package org.tutora.userservice.service.implementation;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.apache.tika.mime.MimeType;
import org.apache.tika.mime.MimeTypeException;
import org.apache.tika.mime.MimeTypes;
import org.tutora.userservice.service.contract.AvatarService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.MediaType;
import org.springframework.stereotype.Service;
import org.springframework.web.reactive.function.client.WebClient;
import reactor.core.publisher.Mono;
import software.amazon.awssdk.core.exception.SdkException;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.presigner.S3Presigner;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.presigner.model.PutObjectPresignRequest;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.time.Duration;
import java.util.HashSet;
import java.util.Set;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class AwsS3AvatarService implements AvatarService {

    private final S3Presigner s3Presigner;
    private final S3Client s3Client;
    private final WebClient webClient;

    @Value("${app.aws.s3.bucket}")
    private String bucketName;

    @Value("${app.aws.cloudfront.domain}")
    private String cloudFrontDomain;

    @Value("${app.avatar.allowed-mime-types}")
    private Set<String> allowedAvatarMimeTypes = new HashSet<>();

    private static final MimeTypes allTypes = MimeTypes.getDefaultMimeTypes();

    @Override
    public String generateUploadUrl(String key, String contentType) {
        PutObjectRequest objectRequest = PutObjectRequest.builder()
                .bucket(bucketName)
                .key(key)
                .contentType(contentType)
                .build();

        PutObjectPresignRequest presignRequest = PutObjectPresignRequest.builder()
                .signatureDuration(Duration.ofMinutes(5))
                .putObjectRequest(objectRequest)
                .build();

        return s3Presigner.presignPutObject(presignRequest).url().toString();
    }

    @Override
    public String getAvatarUrl(String avatarKey) {
        if (avatarKey == null) {
            return null; // lub URL do domyślnego avatara
        }
        return "https://" + cloudFrontDomain + "/" + avatarKey;
    }

    @Override
    public void deleteAvatar(String key) {
        s3Client.deleteObject(b -> b.bucket(bucketName).key(key));
    }

    @Override
    public String saveAvatarFromUrl(UUID userId, String sourceUrl) {
        return downloadAndSaveAvatar(userId, sourceUrl).block();
    }

    private Mono<String> downloadAndSaveAvatar(UUID userId, String sourceUrl) {
        return webClient.get()
                .uri(sourceUrl)
                .exchangeToMono(response -> {
                    if (response.statusCode().isError()) {
                        log.error("Failed to download avatar from {}. Status: {}", sourceUrl, response.statusCode());
                        return Mono.error(new RuntimeException("Błąd podczas pobierania avatara. Status: " + response.statusCode()));
                    }

                    MediaType contentType = response.headers().contentType().orElse(MediaType.APPLICATION_OCTET_STREAM);
                    validateContentType(contentType);

                    String key = generateS3Key(userId, contentType.toString());

                    return response.bodyToMono(byte[].class).map(bytes -> {
                        InputStream inputStream = new ByteArrayInputStream(bytes);
                        return uploadToS3(key, contentType.toString(), bytes.length, inputStream);
                    });
                });
    }

    private void validateContentType(MediaType contentType) {
        if (!allowedAvatarMimeTypes.contains(contentType.toString())) {
            throw new IllegalArgumentException("Niedozwolony typ pliku: " + contentType + ". Dozwolone typy to: " + allowedAvatarMimeTypes);
        }
    }

    private String generateS3Key(UUID userId, String contentType) {
        String extension;
        try {
            MimeType mimeType = allTypes.forName(contentType);
            extension = mimeType.getExtension();
        } catch (MimeTypeException e) {
            log.warn("Could not determine extension for MIME type: {}. Defaulting to empty.", contentType, e);
            extension = "";
        }
        return "avatars/" + userId + "/" + UUID.randomUUID() + extension;
    }

    private String uploadToS3(String key, String contentType, long contentLength, InputStream data) {
        try {
            PutObjectRequest request = PutObjectRequest.builder()
                    .bucket(bucketName)
                    .key(key)
                    .contentType(contentType)
                    .build();

            s3Client.putObject(request, RequestBody.fromInputStream(data, contentLength));
            log.info("Successfully saved avatar to S3 with key: {}", key);
            return key;
        } catch (SdkException e) {
            log.error("Failed to upload avatar to S3 with key: {}", key, e);
            throw new RuntimeException("Nie udało się zapisać avatara w S3", e);
        }
    }
}
