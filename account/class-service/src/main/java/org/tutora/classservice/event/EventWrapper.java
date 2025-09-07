package org.tutora.classservice.event;

import java.io.Serializable;

public record EventWrapper<T extends Event> (
        String pattern,
        T data
) implements Serializable {
}
