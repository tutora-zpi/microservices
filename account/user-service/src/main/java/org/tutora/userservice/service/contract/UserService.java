package org.tutora.userservice.service.contract;

import org.tutora.userservice.dto.UpdateUserDto;
import org.tutora.userservice.entity.User;
import org.tutora.userservice.security.userinfo.OAuth2UserInfo;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;

import java.util.Optional;
import java.util.UUID;

public interface UserService {
    User findById(UUID id);
    Page<User> findByNameStartingWithIgnoreCaseOrSurnameStartingWithIgnoreCase(
            String query,
            Pageable pageable
    );
    Optional<User> findByProviderAndProviderId(String provider, String providerId);
    User registerUser(OAuth2UserInfo userInfo);
    User save(User user);
    String updateUserAvatar(UUID userId, String contentType);
    User updateUserData(UUID userId, UpdateUserDto dto);
}
