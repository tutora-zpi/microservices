package org.example.userservice.dto;

import java.util.Set;

public record UserDto(
        Long id,
        String email,
        Set<String> roles
) {}

