package org.tutora.classservice.dto;

import java.io.Serializable;

public record Event<T> (
        String pattern,
        T data
) implements Serializable {
}
