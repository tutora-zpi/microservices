package org.example.userservice.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.userservice.dto.AvatarUploadRequest;
import org.example.userservice.dto.UpdateUserDto;
import org.example.userservice.dto.UserDto;
import org.example.userservice.entity.User;
import org.example.userservice.mapper.UserMapper;
import org.example.userservice.service.contract.UserService;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.Map;
import java.util.UUID;

@RestController
@RequestMapping("/users")
@RequiredArgsConstructor
@Slf4j
public class UserController {

    private final UserService userService;
    private final UserMapper userMapper;

    @GetMapping("/{id}")
    public ResponseEntity<UserDto> getUser(@PathVariable UUID id) {
        log.info("Request to get user with id {}", id);

        User user = userService.findById(id);

        return ResponseEntity.ok(userMapper.toDto(user));
    }

    @PatchMapping("/{id}")
    public ResponseEntity<UserDto> updateUser(
            @PathVariable UUID id,
            @RequestBody UpdateUserDto userDto) {
        log.info("Request to update user {}", userDto);

        User updated = userService.updateUserData(id, userDto);

        return ResponseEntity.ok(userMapper.toDto(updated));
    }

    @PostMapping("/{id}/avatar")
    public ResponseEntity<Map<String, String>> generateAvatarPresignedUrl(
            @PathVariable UUID id,
            @RequestBody AvatarUploadRequest request) {
        log.info("Request to generate avatar presigned url for user with id {}", id);

        String uploadUrl = userService.updateUserAvatar(id, request.getContentType());

        return ResponseEntity.ok(Map.of("uploadUrl", uploadUrl));
    }
}
