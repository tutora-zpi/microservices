package org.tutora.userservice.event;

import java.io.Serializable;

public record EventWrapper<T extends Event> (
        String pattern,
        T data
) implements Serializable {
}
