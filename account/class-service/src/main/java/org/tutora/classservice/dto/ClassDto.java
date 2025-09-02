package org.tutora.classservice.dto;

import java.util.List;
import java.util.UUID;

public record ClassDto(
        UUID id,
        String name,
        List<UserRoleDto> users
) {
}
