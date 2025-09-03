package org.tutora.classservice.service.contract;

import org.springframework.security.oauth2.jwt.Jwt;

public interface AuthService {
    String getUserId(Jwt principal);
    String getFullName(Jwt principal);
}
