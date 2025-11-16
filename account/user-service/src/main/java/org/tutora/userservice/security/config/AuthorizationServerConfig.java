package org.tutora.userservice.security.config;

import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.RSAKey;
import com.nimbusds.jose.jwk.source.ImmutableJWKSet;
import com.nimbusds.jose.jwk.source.JWKSource;
import com.nimbusds.jose.proc.SecurityContext;
import lombok.RequiredArgsConstructor;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.context.properties.EnableConfigurationProperties;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.core.annotation.Order;
import org.springframework.http.MediaType;
import org.springframework.security.config.Customizer;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.core.AuthorizationGrantType;
import org.springframework.security.oauth2.core.ClientAuthenticationMethod;
import org.springframework.security.oauth2.server.authorization.client.InMemoryRegisteredClientRepository;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClient;
import org.springframework.security.oauth2.server.authorization.client.RegisteredClientRepository;
import org.springframework.security.oauth2.server.authorization.config.annotation.web.configurers.OAuth2AuthorizationServerConfigurer;
import org.springframework.security.oauth2.server.authorization.settings.AuthorizationServerSettings;
import org.springframework.security.oauth2.server.authorization.settings.ClientSettings;
import org.springframework.security.web.SecurityFilterChain;
import org.springframework.security.web.authentication.LoginUrlAuthenticationEntryPoint;
import org.springframework.security.web.util.matcher.MediaTypeRequestMatcher;
import org.tutora.userservice.security.config.properties.BotClientProperties;
import org.tutora.userservice.security.config.properties.WebClientProperties;
import org.tutora.userservice.security.key.RsaKeyProperties;

import java.util.UUID;

/**
 * Konfiguracja łańcucha filtrów dla Serwera Autoryzacji (OAuth2).
 * Ten łańcuch ma najwyższy priorytet (@Order(1)) i obsługuje
 * wyłącznie endpointy związane z protokołem OAuth2 (np. /oauth2/token).
 */
@Configuration
@Order(1)
@RequiredArgsConstructor
@EnableConfigurationProperties({WebClientProperties.class, BotClientProperties.class})
public class AuthorizationServerConfig {

    private final WebClientProperties webClientProperties;
    private final BotClientProperties botClientProperties;

    @Value("${app.oauth2.issuer-uri}")
    private String issuerUri;

    @Bean
    public SecurityFilterChain authorizationServerSecurityFilterChain(HttpSecurity http) throws Exception {

        http.with(OAuth2AuthorizationServerConfigurer.authorizationServer(), configurer ->
                configurer.oidc(Customizer.withDefaults()));

        http
                .exceptionHandling(exceptions ->
                        exceptions.defaultAuthenticationEntryPointFor(
                                new LoginUrlAuthenticationEntryPoint("/login"),
                                new MediaTypeRequestMatcher(MediaType.TEXT_HTML)
                        )
                )
                .cors(Customizer.withDefaults())
                .formLogin(Customizer.withDefaults());

        return http.build();
    }

    @Bean
    public AuthorizationServerSettings authorizationServerSettings() {
        return AuthorizationServerSettings.builder()
                .issuer(issuerUri)
                .jwkSetEndpoint("/.well-known/jwks.json")
                .build();
    }

    @Bean
    public PasswordEncoder passwordEncoder() {
        return new BCryptPasswordEncoder();
    }

    @Bean
    public RegisteredClientRepository registeredClientRepository(PasswordEncoder passwordEncoder) {
        RegisteredClient webClient = RegisteredClient
                .withId(UUID.randomUUID().toString())
                .clientId(webClientProperties.clientId())
                .clientAuthenticationMethod(ClientAuthenticationMethod.NONE)
                .authorizationGrantType(AuthorizationGrantType.AUTHORIZATION_CODE)
                .authorizationGrantType(AuthorizationGrantType.REFRESH_TOKEN)
                .redirectUris(properties -> properties.addAll(webClientProperties.redirectUris()))
                .scopes(scopes -> scopes.addAll(webClientProperties.scopes()))
                .clientSettings(ClientSettings.builder().requireProofKey(true).build())
                .build();

        RegisteredClient recordingBot = RegisteredClient
                .withId(UUID.randomUUID().toString())
                .clientId(botClientProperties.clientId())
                .clientSecret(passwordEncoder.encode(botClientProperties.clientSecret()))
                .clientAuthenticationMethod(ClientAuthenticationMethod.CLIENT_SECRET_BASIC)
                .authorizationGrantType(AuthorizationGrantType.CLIENT_CREDENTIALS)
                .scopes(scopes -> scopes.addAll(botClientProperties.scopes()))
                .build();

        return new InMemoryRegisteredClientRepository(webClient, recordingBot);
    }

    @Bean
    public JWKSource<SecurityContext> jwkSource(RsaKeyProperties rsaKeyProperties) {
        RSAKey rsaKey = new RSAKey.Builder(rsaKeyProperties.publicKey())
                .privateKey(rsaKeyProperties.privateKey())
                .keyID(UUID.randomUUID().toString())
                .build();

        JWKSet jwkSet = new JWKSet(rsaKey);
        return new ImmutableJWKSet<>(jwkSet);
    }
}
