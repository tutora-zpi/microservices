package org.tutora.userservice.event;

import java.util.List;

public record UserDetailsRespondedEvent(
        String notificationId,
        UserDetails sender,
        UserDetails receiver
) implements Event {

    @Override
    public String name() {
        return this.getClass().getSimpleName();
    }

    public record UserDetails(
            String id,
            String firstName,
            String lastName,
            List<String> roles
    ) {}
}
