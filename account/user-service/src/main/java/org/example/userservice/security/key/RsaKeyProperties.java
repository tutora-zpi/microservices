package org.example.userservice.security.key;

import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;

public record RsaKeyProperties(RSAPublicKey publicKey, RSAPrivateKey privateKey) {
}

