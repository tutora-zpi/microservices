package org.tutora.userservice.event;

public record UserDetailsRequestedEvent(
    String notificationId,
    String senderId,
    String receiverId
) implements Event {
}
