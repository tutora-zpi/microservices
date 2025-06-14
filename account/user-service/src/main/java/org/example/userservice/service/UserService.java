package org.example.userservice.service;

import org.example.userservice.entity.User;

import java.util.Optional;

public interface UserService {
    User findById(Long id);
    User findByEmail(String email);
    Optional<User> findByProviderAndProviderId(String provider, String providerId);
    User save(User user);
}
