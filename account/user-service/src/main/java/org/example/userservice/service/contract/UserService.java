package org.example.userservice.service.contract;

import org.example.userservice.dto.UpdateUserDto;
import org.example.userservice.entity.User;
import org.example.userservice.security.userinfo.OAuth2UserInfo;

import java.util.Optional;
import java.util.UUID;

public interface UserService {
    User findById(UUID id);
    User findByNameAndSurname(String name, String surname);
    Optional<User> findByProviderAndProviderId(String provider, String providerId);
    User registerUser(OAuth2UserInfo userInfo);
    User save(User user);
    String updateUserAvatar(UUID userId, String contentType);
    User updateUserData(UUID userId, UpdateUserDto dto);
}
