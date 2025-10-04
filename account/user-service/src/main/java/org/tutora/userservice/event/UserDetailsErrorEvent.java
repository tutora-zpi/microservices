package org.tutora.userservice.event;

public record UserDetailsErrorEvent(
        String notificationId,
        String errorMessage
) implements Event {
}