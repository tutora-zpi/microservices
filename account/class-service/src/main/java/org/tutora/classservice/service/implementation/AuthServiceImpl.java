package org.tutora.classservice.service.implementation;

import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.stereotype.Service;
import org.tutora.classservice.service.contract.AuthService;

@Service
public class AuthServiceImpl implements AuthService {
    @Override
    public String getUserId(Jwt principal) {
        return principal.getSubject();
    }
}
