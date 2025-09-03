package org.tutora.classservice.dto;

import java.util.UUID;

public record ClassInvitationEvent(
        String senderFullName,
        String className,
        UUID receiverId
) {
}
