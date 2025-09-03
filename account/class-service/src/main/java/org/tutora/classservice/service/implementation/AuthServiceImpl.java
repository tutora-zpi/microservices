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

    @Override
    public String getFullName(Jwt principal) {
        String firstName = principal.getClaim("first_name");
        String lastName = principal.getClaim("last_name");

        return firstName + " " + lastName;
    }
}
