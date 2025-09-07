package org.tutora.classservice.dto;

import java.time.LocalDateTime;

public record ErrorDetailsDto(
        LocalDateTime timestamp,
        String message,
        String path,
        int status
) {}