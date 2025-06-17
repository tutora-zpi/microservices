package org.example.userservice.dto;

import java.util.Set;
import java.util.UUID;

public record UserDto(
        UUID id,
        String email,
        Set<String> roles
) {}

