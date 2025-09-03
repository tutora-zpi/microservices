package org.tutora.classservice.controller;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;
import org.tutora.classservice.dto.InvitationDto;
import org.tutora.classservice.entity.Invitation;
import org.tutora.classservice.mapper.InvitationMapper;
import org.tutora.classservice.service.contract.AuthService;
import org.tutora.classservice.service.contract.InvitationService;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/invitations")
@RequiredArgsConstructor
@Slf4j
public class InvitationController {

    private final InvitationService invitationService;
    private final AuthService authService;
    private final InvitationMapper invitationMapper;

    @GetMapping("/me")
    public ResponseEntity<List<InvitationDto>> getMyInvitations(@AuthenticationPrincipal Jwt principal) {
        UUID userId = UUID.fromString(authService.getUserId(principal));
        log.info("Request for fetching invitations for current user [{}]", userId);

        List<InvitationDto> invitationDtos = invitationService.getInvitationsForUser(userId)
                .stream().map(invitationMapper::toDto).toList();

        return ResponseEntity
                .ok(invitationDtos);
    }

    @GetMapping("/classes/{classId}")
    public ResponseEntity<List<InvitationDto>> getClassInvitations(@PathVariable UUID classId) {
        log.info("Request for fetching invitations for class [{}]", classId);

        List<InvitationDto> invitationDtos = invitationService.getInvitationsForClass(classId)
                .stream().map(invitationMapper::toDto).toList();

        return ResponseEntity
                .ok(invitationDtos);
    }

    @PostMapping("/{classId}/users/{userId}")
    public ResponseEntity<InvitationDto> inviteUser(
            @AuthenticationPrincipal Jwt principal,
            @PathVariable UUID classId,
            @PathVariable UUID userId
    ) {
        String senderFullName = authService.getFullName(principal);
        log.info("User [{}] invites [{}] to class [{}]", senderFullName, userId, classId);

        Invitation inv = invitationService.inviteUser(senderFullName, classId, userId);
        return ResponseEntity
                .status(HttpStatus.CREATED)
                .body(invitationMapper.toDto(inv));
    }

    @DeleteMapping("/{classId}/users/{userId}")
    public ResponseEntity<Void> cancelInvitation(
            @PathVariable UUID classId,
            @PathVariable UUID userId
    ) {
        log.info("Cancel invitation for user [{}] in class [{}]", userId, classId);

        invitationService.cancelInvitation(classId, userId);
        return ResponseEntity
                .noContent()
                .build();
    }

    @PostMapping("/{classId}/users/{userId}/accept")
    public ResponseEntity<Void> acceptInvitation(
            @PathVariable UUID classId,
            @PathVariable UUID userId
    ) {
        log.info("User [{}] accepts invitation to class [{}]", userId, classId);

        invitationService.joinClass(classId, userId);
        return ResponseEntity.ok().build();
    }

    @PostMapping("/{classId}/users/{userId}/decline")
    public ResponseEntity<Void> declineInvitation(
            @PathVariable UUID classId,
            @PathVariable UUID userId
    ) {
        log.info("User [{}] declines invitation to class [{}]", userId, classId);

        invitationService.declineInvitation(classId, userId);
        return ResponseEntity.ok().build();
    }
}
