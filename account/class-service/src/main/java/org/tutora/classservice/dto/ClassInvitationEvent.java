package org.tutora.classservice.dto;

import java.util.UUID;

public record ClassInvitationEvent(
        UUID senderId,
        String className,
        UUID receiverId
) {
}
