package org.tutora.classservice.dto;

import java.util.UUID;

public record ClassCreateRequest(
        String name,
        UUID hostId
) {
}
