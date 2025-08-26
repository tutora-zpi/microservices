package org.example.userservice.service.contract;

import org.example.userservice.entity.User;

import java.util.Optional;
import java.util.UUID;

public interface UserService {
    User findById(UUID id);
    User findByEmail(String email);
    Optional<User> findByProviderAndProviderId(String provider, String providerId);
    User save(User user);
    String updateUserAvatar(User user, String contentType);
}
