package org.tutora.classservice.exception;

import java.util.UUID;

public class UserRejectedInvitationException extends RuntimeException {
    public UserRejectedInvitationException(UUID userId, UUID classId) {
        super("User " + userId + " has already rejected invitation to class " + classId);
    }
}
