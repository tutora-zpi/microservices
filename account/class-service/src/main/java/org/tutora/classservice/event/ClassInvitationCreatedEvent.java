package org.tutora.classservice.event;

import org.tutora.classservice.dto.UserDto;

public record ClassInvitationCreatedEvent(
        String className,
        UserDto receiver,
        UserDto sender
) implements Event {
    @Override
    public String name() {
        return this.getClass().getSimpleName();
    }
}
