package org.tutora.classservice.exception;

public class UnauthorizedActionException extends RuntimeException {
    public UnauthorizedActionException(String resource, Object resourceId, String action) {
        super(String.format("User is not authorized to %s on %s with id '%s'", action, resource, resourceId));
    }
}
