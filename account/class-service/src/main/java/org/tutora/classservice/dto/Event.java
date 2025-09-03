package org.tutora.classservice.dto;

public record Event<T>(
        String pattern,
        T data
) {
}
