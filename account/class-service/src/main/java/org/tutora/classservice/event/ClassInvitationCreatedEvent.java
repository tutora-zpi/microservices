package org.tutora.classservice.event;

public record ClassInvitationCreatedEvent(
        String senderId,
        String className,
        String receiverId
) implements Event {
    @Override
    public String name() {
        return this.getClass().getSimpleName();
    }
}
