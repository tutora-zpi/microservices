package org.tutora.classservice.dto;

import java.util.UUID;

public record UserRoleDto(
        UUID userId,
        String role
) {
}
