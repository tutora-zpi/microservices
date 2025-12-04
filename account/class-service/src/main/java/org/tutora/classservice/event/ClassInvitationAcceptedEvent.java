package org.tutora.classservice.event;

import java.util.UUID;

public record ClassInvitationAcceptedEvent(
        UUID classId,
        String className,
        UUID accepterId,
        UUID roomHostId
) implements Event {
    @Override
    public String name() {
        return this.getClass().getSimpleName();
    }
}
