package org.example.userservice.service;

import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;
import software.amazon.awssdk.services.s3.S3Client;
import software.amazon.awssdk.services.s3.presigner.S3Presigner;
import software.amazon.awssdk.services.s3.model.PutObjectRequest;
import software.amazon.awssdk.services.s3.presigner.model.PresignedPutObjectRequest;
import software.amazon.awssdk.services.s3.presigner.model.PutObjectPresignRequest;

import java.time.Duration;

@Service
@RequiredArgsConstructor
public class AwsS3AvatarService implements AvatarService {

    private final S3Presigner s3Presigner;
    private final S3Client s3Client;

    @Value("${app.s3.bucket}")
    private String bucketName;

    @Value("${app.cloudfront.domain}")
    private String cloudFrontDomain;

    @Override
    public String generateUploadUrl(Long userId, String contentType) {
        String key = "avatars/" + userId + ".png";

        PutObjectRequest objectRequest = PutObjectRequest.builder()
                .bucket(bucketName)
                .key(key)
                .contentType(contentType)
                .build();

        PutObjectPresignRequest presignRequest = PutObjectPresignRequest.builder()
                .signatureDuration(Duration.ofMinutes(5))
                .putObjectRequest(objectRequest)
                .build();

        PresignedPutObjectRequest presignedRequest = s3Presigner.presignPutObject(presignRequest);

        return presignedRequest.url().toString();
    }

    @Override
    public String getAvatarUrl(String avatarKey) {
        if (avatarKey == null) {
            return null; // albo "https://cdn.myapp.com/default-avatar.png"
        }
        return "https://" + cloudFrontDomain + "/" + avatarKey;
    }
}
