package org.example.userservice.exception;

import org.springframework.security.core.AuthenticationException;

/**
 * Wyjątek rzucany w przypadku błędu podczas przetwarzania danych
 * użytkownika otrzymanych od dostawcy OAuth2.
 *
 * Przykład: Dostawca nie zwrócił adresu e-mail, który jest wymagany
 * do założenia konta w naszym systemie.
 */
public class OAuth2AuthenticationProcessingException extends AuthenticationException {

    public OAuth2AuthenticationProcessingException(String msg) {
        super(msg);
    }

    public OAuth2AuthenticationProcessingException(String msg, Throwable cause) {
        super(msg, cause);
    }
}