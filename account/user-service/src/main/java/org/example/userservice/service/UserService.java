package org.example.userservice.service;

import org.example.userservice.entity.User;

import java.util.Optional;
import java.util.UUID;

public interface UserService {
    User findById(UUID id);
    User findByEmail(String email);
    Optional<User> findByProviderAndProviderId(String provider, String providerId);
    User save(User user);
}
