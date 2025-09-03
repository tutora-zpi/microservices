package org.tutora.classservice.dto;

import java.time.LocalDateTime;
import java.util.UUID;

public record InvitationDto(
        UUID classId,
        UUID userId,
        String status,
        LocalDateTime createdAt
) {
}
