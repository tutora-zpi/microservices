package org.tutora.userservice.security.key;

import org.springframework.boot.context.properties.ConfigurationProperties;

@ConfigurationProperties(prefix = "app.rsa")
public record RsaKeyPaths(String publicKeyPath, String privateKeyPath) {
}
