package org.tutora.classservice.exception;

import java.util.UUID;

public class UserAlreadyInvitedException extends RuntimeException {
    public UserAlreadyInvitedException(UUID userId, UUID classId) {
        super("User " + userId + " has already been invited to class " + classId);
    }
}
