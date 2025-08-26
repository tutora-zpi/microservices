package org.example.userservice.controller;

import lombok.RequiredArgsConstructor;
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
public class UserController {

    private final UserService userService;
    private final UserMapper userMapper;

    @GetMapping("/{id}")
    public ResponseEntity<UserDto> getUser(@PathVariable UUID id) {
        User user = userService.findById(id);

        return ResponseEntity.ok(userMapper.toDto(user));
    }

    @PatchMapping("/{id}")
    public ResponseEntity<UserDto> updateUser(
            @PathVariable UUID id,
            @RequestBody UpdateUserDto userDto) {
        User updated = userService.updateUserData(id, userDto);

        return ResponseEntity.ok(userMapper.toDto(updated));
    }

    @PostMapping("/{id}/avatar")
    public ResponseEntity<Map<String, String>> updateAvatar(
            @PathVariable UUID id,
            @RequestBody AvatarUploadRequest request) {
        String uploadUrl = userService.updateUserAvatar(id, request.getContentType());

        return ResponseEntity.ok(Map.of("uploadUrl", uploadUrl));
    }
}
