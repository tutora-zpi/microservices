package org.example.userservice.config;

import org.example.userservice.security.RsaKeyProperties;
import org.example.userservice.security.RsaKeyPaths;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import java.nio.file.Files;
import java.nio.file.Paths;
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
        var publicKeyBytes = Files.readAllBytes(Paths.get(ClassLoader.getSystemResource(rsaKeyPaths.publicKeyPath().replace("classpath:", "")).toURI()));
        var privateKeyBytes = Files.readAllBytes(Paths.get(ClassLoader.getSystemResource(rsaKeyPaths.privateKeyPath().replace("classpath:", "")).toURI()));

        var publicKeySpec = new X509EncodedKeySpec(stripHeadersAndDecode(new String(publicKeyBytes)));
        var privateKeySpec = new PKCS8EncodedKeySpec(stripHeadersAndDecode(new String(privateKeyBytes)));

        KeyFactory keyFactory = KeyFactory.getInstance("RSA");

        RSAPublicKey publicKey = (RSAPublicKey) keyFactory.generatePublic(publicKeySpec);
        RSAPrivateKey privateKey = (RSAPrivateKey) keyFactory.generatePrivate(privateKeySpec);

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