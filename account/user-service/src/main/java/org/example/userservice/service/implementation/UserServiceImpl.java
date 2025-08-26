package org.example.userservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.example.userservice.dto.UpdateUserDto;
import org.example.userservice.entity.User;
import org.example.userservice.exception.ResourceNotFoundException; // Zmieniony import
import org.example.userservice.mapper.UserMapper;
import org.example.userservice.repository.UserRepository;
import org.example.userservice.service.contract.AvatarService;
import org.example.userservice.service.contract.UserService;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Optional;
import java.util.UUID;

@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class UserServiceImpl implements UserService {

    private final UserRepository userRepository;
    private final AvatarService avatarService;
    private final UserMapper userMapper;

    @Override
    public User findById(UUID id) {
        return userRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("User", "id", id));
    }

    @Override
    public User findByNameAndSurname(String name, String surname) {
        return userRepository.findByNameAndSurname(name, surname)
                .orElseThrow(() -> new ResourceNotFoundException("User", "name and surname", name + surname));
    }

    @Override
    public Optional<User> findByProviderAndProviderId(String provider, String providerId) {
        return userRepository.findByProviderAndProviderId(provider, providerId);
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

        String newKey = "avatars/" + user.getId() + "/" + UUID.randomUUID() + ".png";

        String uploadUrl = avatarService.generateUploadUrl(newKey, contentType);

        user.setAvatarKey(newKey);
        userRepository.save(user);

        return uploadUrl;
    }

    @Override
    public User updateUserData(UUID userId, UpdateUserDto dto) {
        User user = userRepository.findById(userId)
                .orElseThrow(() -> new ResourceNotFoundException("User", "id", userId));

        userMapper.updateUserFromDto(dto, user);

        return userRepository.save(user);
    }
}