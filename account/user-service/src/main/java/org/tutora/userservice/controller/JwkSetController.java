package org.tutora.userservice.controller;

import com.nimbusds.jose.jwk.JWKSet;
import com.nimbusds.jose.jwk.RSAKey;
import lombok.RequiredArgsConstructor;
import org.tutora.userservice.security.key.RsaKeyProperties;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
@RequiredArgsConstructor
public class JwkSetController {

    private final RsaKeyProperties rsaKeyProperties;

    /**
     * Publiczny endpoint udostępniający klucz publiczny RSA w formacie JWKS.
     * Inne mikroserwisy będą pobierać ten klucz, aby weryfikować podpisy tokenów JWT.
     */
    @GetMapping("/.well-known/jwks.json")
    public Map<String, Object> keys() {
        RSAKey key = new RSAKey.Builder(rsaKeyProperties.publicKey())
                .keyID("rsa-key-1")
                .build();

        JWKSet jwkSet = new JWKSet(key);

        return jwkSet.toJSONObject();
    }
}