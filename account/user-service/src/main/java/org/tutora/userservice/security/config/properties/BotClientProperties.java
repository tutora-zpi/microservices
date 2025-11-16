package org.tutora.userservice.security.config.properties;

import org.springframework.boot.context.properties.ConfigurationProperties;
import java.util.List;

@ConfigurationProperties(prefix = "app.oauth2.clients.recording-bot")
public record BotClientProperties(
        String clientId,
        String clientSecret,
        List<String> scopes
) {}
