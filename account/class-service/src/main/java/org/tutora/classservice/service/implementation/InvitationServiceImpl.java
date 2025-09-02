package org.tutora.classservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.tutora.classservice.entity.*;
import org.tutora.classservice.exception.ResourceNotFoundException;
import org.tutora.classservice.repository.InvitationRepository;
import org.tutora.classservice.repository.InvitationStatusRepository;
import org.tutora.classservice.service.contract.ClassService;
import org.tutora.classservice.service.contract.InvitationService;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class InvitationServiceImpl implements InvitationService {

    private final InvitationRepository invitationRepository;
    private final InvitationStatusRepository invitationStatusRepository;

    private final ClassService classService;

    @Override
    public void inviteUser(UUID classId, UUID userId) {
        saveInvitation(classId, userId);

        //TODO: emit event to notification service
    }

    @Override
    public void cancelInvitation(UUID classId, UUID userId) {
        invitationRepository.deleteByClassroomIdAndUserId(classId, userId);

        //TODO: emit event to notification service
    }

    @Override
    public void joinClass(UUID classId, UUID userId) {
        Invitation inv = getInvitation(classId, userId);

        inv.setStatus(getInvitationStatus(InvitationStatusName.ACCEPTED));
        invitationRepository.save(inv);

        classService.addUserToClass(classId, userId, RoleName.GUEST);
    }

    @Override
    public void declineInvitation(UUID classId, UUID userId) {
        Invitation inv = getInvitation(classId, userId);

        inv.setStatus(getInvitationStatus(InvitationStatusName.DECLINED));
        invitationRepository.save(inv);

        //TODO: emit event to notification service
    }

    @Override
    public List<Invitation> getInvitationsForUser(UUID userId) {
        return invitationRepository
                .findAllByUserIdAndStatus(userId, getInvitationStatus(InvitationStatusName.INVITED));
    }

    @Override
    public List<Invitation> getInvitationsForClass(UUID classId) {
        return invitationRepository
                .findAllByClassroomIdAndStatus(classId, getInvitationStatus(InvitationStatusName.INVITED));
    }


    private InvitationStatus getInvitationStatus(InvitationStatusName statusName) {
        return invitationStatusRepository.findByStatusName(statusName)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "Invitation status", "status name", statusName
                ));
    }

    private void saveInvitation(UUID classId, UUID userId) {
        InvitationStatus status = getInvitationStatus(InvitationStatusName.INVITED);

        Invitation inv = Invitation.builder()
                .classroom(classService.getClassById(classId))
                .userId(userId)
                .status(status)
                .createdAt(LocalDateTime.now())
                .build();

        invitationRepository.save(inv);
    }

    private Invitation getInvitation(UUID classId, UUID userId) {
        return invitationRepository.findByClassroomIdAndUserId(classId, userId)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "Invitation", "classId" + classId.toString() + " and userId", userId));
        //TODO refactor exception constructor to accept many objects
    }
}
