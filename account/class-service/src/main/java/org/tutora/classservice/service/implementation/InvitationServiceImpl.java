package org.tutora.classservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.tutora.classservice.client.NotificationPublisher;
import org.tutora.classservice.dto.ClassInvitationEvent;
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
    private final NotificationPublisher notificationPublisher;

    private final ClassService classService;

    @Override
    public Invitation inviteUser(String senderFullName, UUID classId, UUID userId) {
        Invitation inv = saveInvitation(classId, userId);

        Classroom classroom = classService.getClassById(classId);

        notificationPublisher.sendClassInvitation(new ClassInvitationEvent(
                senderFullName,
                classroom.getName(),
                userId
        ));

        return inv;
    }

    @Override
    public void cancelInvitation(UUID classId, UUID userId) {
        invitationRepository.deleteByClassroomIdAndUserId(classId, userId);
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

    private Invitation saveInvitation(UUID classId, UUID userId) {
        InvitationStatus status = getInvitationStatus(InvitationStatusName.INVITED);

        Invitation inv = Invitation.builder()
                .classroom(classService.getClassById(classId))
                .userId(userId)
                .status(status)
                .createdAt(LocalDateTime.now())
                .build();

        return invitationRepository.save(inv);
    }

    private Invitation getInvitation(UUID classId, UUID userId) {
        return invitationRepository.findByClassroomIdAndUserId(classId, userId)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "Invitation", "classId" + classId.toString() + " and userId", userId));
        //TODO refactor exception constructor to accept multiple objects
    }
}
