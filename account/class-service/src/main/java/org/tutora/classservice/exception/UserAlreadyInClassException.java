package org.tutora.classservice.exception;

import java.util.UUID;

public class UserAlreadyInClassException extends RuntimeException {
    public UserAlreadyInClassException(UUID userId, UUID classId) {
        super("User " + userId + " is already a member of class " + classId);
    }
}
