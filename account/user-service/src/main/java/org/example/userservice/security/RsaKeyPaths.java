package org.example.userservice.security;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "app.rsa")
public record RsaKeyPaths(String publicKeyPath, String privateKeyPath) {
}
