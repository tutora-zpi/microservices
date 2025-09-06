package org.example.userservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.example.userservice.dto.UpdateUserDto;
import org.example.userservice.entity.Role;
import org.example.userservice.entity.RoleName;
import org.example.userservice.entity.User;
import org.example.userservice.exception.ResourceNotFoundException;
import org.example.userservice.repository.RoleRepository;
import org.example.userservice.repository.UserRepository;
import org.example.userservice.security.userinfo.OAuth2UserInfo;
import org.example.userservice.service.contract.AvatarService;
import org.example.userservice.service.contract.UserService;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;
import java.util.Set;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final RoleRepository roleRepository;
    private final AvatarService avatarService;

    @Override
    public User findById(UUID id) {
        return userRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("User", "id", id));
    }

    @Override
    public Page<User> findByNameStartingWithIgnoreCaseOrSurnameStartingWithIgnoreCase(
            String query,
            Pageable pageable) {
        return userRepository.findByNameStartingWithIgnoreCaseOrSurnameStartingWithIgnoreCase(
                query,
                query,
                pageable);
    }

    @Override
    public Optional<User> findByProviderAndProviderId(String provider, String providerId) {
        return userRepository.findByProviderAndProviderId(provider, providerId);
    }

    @Override
    public User registerUser(OAuth2UserInfo userInfo) {
        User newUser = new User();
        newUser.setId(UUID.randomUUID());
        newUser.setProvider(userInfo.getProvider());
        newUser.setProviderId(userInfo.getId());
        newUser.setEmail(userInfo.getEmail());
        newUser.setName(userInfo.getName());
        newUser.setSurname(userInfo.getSurname());

        Role userRole = roleRepository.findByName(RoleName.USER)
                .orElseThrow(() -> new ResourceNotFoundException("Role", "name", RoleName.USER));
        newUser.setRoles(Set.of(userRole));

        String avatarKey = avatarService.saveAvatarFromUrl(newUser.getId(), userInfo.getImageUrl());
        newUser.setAvatarKey(avatarKey);

        return save(newUser);
    }

    @Override
    @Transactional
    public User save(User user) {
        return userRepository.save(user);
    }

    @Override
    public String updateUserAvatar(UUID userId, String contentType) {
        User user = findById(userId);

        if (user.getAvatarKey() != null) {
            avatarService.deleteAvatar(user.getAvatarKey());
        }

        String newKey = "avatars/" + user.getId().toString() + "/" + UUID.randomUUID() + ".png";

        String uploadUrl = avatarService.generateUploadUrl(newKey, contentType);

        user.setAvatarKey(newKey);
        userRepository.save(user);

        return uploadUrl;
    }

    @Override
    @Transactional
    public User updateUserData(UUID userId, UpdateUserDto dto) {
        User user = userRepository.findById(userId)
                .orElseThrow(() -> new ResourceNotFoundException("User", "id", userId));

        user.setEmail(dto.getEmail() == null ? user.getEmail() : dto.getEmail());
        user.setName(dto.getName() == null ? user.getName() : dto.getName());
        user.setSurname(dto.getSurname() == null ? user.getSurname() : dto.getSurname());

        return userRepository.save(user);
    }
}