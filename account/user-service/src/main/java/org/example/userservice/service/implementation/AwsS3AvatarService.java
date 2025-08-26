package org.example.userservice.service.implementation;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.userservice.service.contract.AvatarService;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import org.springframework.web.client.RestTemplate;
import software.amazon.awssdk.core.sync.RequestBody;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.presigner.S3Presigner;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.presigner.model.PutObjectPresignRequest;

import java.time.Duration;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Slf4j
public class AwsS3AvatarService implements AvatarService {

    private final S3Presigner s3Presigner;
    private final S3Client s3Client;
    private final RestTemplate restTemplate;

    @Value("${app.s3.bucket}")
    private String bucketName;

    @Value("${app.cloudfront.domain}")
    private String cloudFrontDomain;

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
            return null; // albo "https://cdn.myapp.com/default-avatar.png"
        }
        return "https://" + cloudFrontDomain + "/" + avatarKey;
    }

    @Override
    public void deleteAvatar(String key) {
        s3Client.deleteObject(b -> b.bucket(bucketName).key(key));
    }

    @Override
    public String saveAvatarFromUrl(UUID userId, String sourceUrl) {
        try {
            byte[] avatarBytes = restTemplate.getForObject(sourceUrl, byte[].class);
            if (avatarBytes == null || avatarBytes.length == 0) {
                throw new RuntimeException("Pusty plik avataru");
            }

            String key = "avatars/" + userId + "/" + UUID.randomUUID() + ".png";

            s3Client.putObject(
                    PutObjectRequest.builder()
                            .bucket(bucketName)
                            .key(key)
                            .contentType("image/jpeg")
                            .build(),
                    RequestBody.fromBytes(avatarBytes)
            );

            return key;
        } catch (Exception e) {
            throw new RuntimeException("Nie udało się pobrać lub zapisać avatara z URL: " + sourceUrl, e);
        }
    }
}
