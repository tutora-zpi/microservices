package org.tutora.classservice.dto;

public record InvitationCreateRequest(
        String classId,
        UserDto receiver,
        UserDto sender
) {
}
