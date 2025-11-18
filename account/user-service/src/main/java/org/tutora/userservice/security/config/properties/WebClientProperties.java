package org.tutora.userservice.security.config.properties;

import org.springframework.boot.context.properties.ConfigurationProperties;
import java.util.List;

@ConfigurationProperties(prefix = "app.oauth2.clients.frontend")
public record WebClientProperties(
        String clientId,
        List<String> redirectUris,
        List<String> scopes
) {}
