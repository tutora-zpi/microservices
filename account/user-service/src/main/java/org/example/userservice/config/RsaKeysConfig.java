package org.example.userservice.config;

import org.example.userservice.security.key.RsaKeyProperties;
import org.example.userservice.security.key.RsaKeyPaths;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.io.ClassPathResource;
import org.springframework.core.io.Resource;

import java.security.KeyFactory;
import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;

@Configuration
@EnableConfigurationProperties(RsaKeyPaths.class)
public class RsaKeysConfig {

    @Bean
    public RsaKeyProperties rsaKeyProperties(RsaKeyPaths rsaKeyPaths) throws Exception {
        Resource privateKeyResource = new ClassPathResource(rsaKeyPaths.privateKeyPath().replace("classpath:", ""));
        Resource publicKeyResource = new ClassPathResource(rsaKeyPaths.publicKeyPath().replace("classpath:", ""));

        byte[] privateKeyBytes = privateKeyResource.getInputStream().readAllBytes();
        byte[] publicKeyBytes = publicKeyResource.getInputStream().readAllBytes();

        String privateKeyPem = new String(privateKeyBytes);
        String publicKeyPem = new String(publicKeyBytes);

        KeyFactory keyFactory = KeyFactory.getInstance("RSA");

        PKCS8EncodedKeySpec privateKeySpec = new PKCS8EncodedKeySpec(stripHeadersAndDecode(privateKeyPem));
        X509EncodedKeySpec publicKeySpec = new X509EncodedKeySpec(stripHeadersAndDecode(publicKeyPem));

        RSAPrivateKey privateKey = (RSAPrivateKey) keyFactory.generatePrivate(privateKeySpec);
        RSAPublicKey publicKey = (RSAPublicKey) keyFactory.generatePublic(publicKeySpec);

        return new RsaKeyProperties(publicKey, privateKey);
    }

    private byte[] stripHeadersAndDecode(String pem) {
        String base64 = pem
                .replaceAll("-----BEGIN (.*)-----", "")
                .replaceAll("-----END (.*)-----", "")
                .replaceAll("\\s", "");
        return Base64.getDecoder().decode(base64);
    }
}