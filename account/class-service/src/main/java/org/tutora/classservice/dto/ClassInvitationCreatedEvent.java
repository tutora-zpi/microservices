package org.tutora.classservice.dto;

import java.util.UUID;

public record ClassInvitationCreatedEvent(
        UUID senderId,
        String className,
        UUID receiverId
) {
}
