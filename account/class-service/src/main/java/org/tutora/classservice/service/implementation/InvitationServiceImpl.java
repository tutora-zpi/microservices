package org.tutora.classservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.tutora.classservice.client.NotificationPublisher;
import org.tutora.classservice.dto.ClassInvitationEvent;
import org.tutora.classservice.entity.*;
import org.tutora.classservice.exception.ResourceNotFoundException;
import org.tutora.classservice.exception.UserAlreadyInClassException;
import org.tutora.classservice.exception.UserAlreadyInvitedException;
import org.tutora.classservice.exception.UserRejectedInvitationException;
import org.tutora.classservice.repository.InvitationRepository;
import org.tutora.classservice.repository.InvitationStatusRepository;
import org.tutora.classservice.repository.UserClassRepository;
import org.tutora.classservice.service.contract.ClassService;
import org.tutora.classservice.service.contract.InvitationService;

import java.time.LocalDateTime;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class InvitationServiceImpl implements InvitationService {

    private final InvitationRepository invitationRepository;
    private final InvitationStatusRepository invitationStatusRepository;
    private final NotificationPublisher notificationPublisher;

    private final ClassService classService;
    private final UserClassRepository userClassRepository;

    @Override
    public Invitation inviteUser(UUID senderId, UUID classId, UUID userId) {
        Classroom classroom = classService.getClassById(classId);

        if (userClassRepository.existsByClassroomIdAndUserId(classId, userId)) {
            throw new UserAlreadyInClassException(userId, classId);
        }

        Optional<Invitation> existingInvitation = invitationRepository.findByClassroomIdAndUserId(classId, userId);
        if(existingInvitation.isPresent()) {
            Invitation inv = existingInvitation.get();
            switch (inv.getStatus().getStatusName()) {
                case DECLINED -> throw new UserRejectedInvitationException(userId, classId);
                case INVITED -> throw new UserAlreadyInvitedException(userId, classId);
            }
        }

        Invitation inv = saveInvitation(classId, userId);

        notificationPublisher.sendClassInvitation(new ClassInvitationEvent(
                senderId,
                classroom.getName(),
                userId
        ));

        return inv;
    }

    @Transactional
    @Override
    public void cancelInvitation(UUID classId, UUID userId) {
        invitationRepository.deleteByClassroomIdAndUserId(classId, userId);
    }

    @Override
    public void joinClass(UUID classId, UUID userId) {
        Invitation inv = getInvitation(classId, userId);

        validateInvitationStatus(inv.getStatus(), classId, userId);

        inv.setStatus(getInvitationStatus(InvitationStatusName.ACCEPTED));
        invitationRepository.save(inv);

        classService.addUserToClass(classId, userId, RoleName.GUEST);
    }

    @Override
    public void declineInvitation(UUID classId, UUID userId) {
        Invitation inv = getInvitation(classId, userId);

        validateInvitationStatus(inv.getStatus(), classId, userId);

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

    private void validateInvitationStatus(InvitationStatus invitationStatus, UUID classId, UUID userId) {
        if (invitationStatus.getStatusName() == InvitationStatusName.ACCEPTED) {
            throw new UserAlreadyInClassException(userId, classId);
        }
        if (invitationStatus.getStatusName() == InvitationStatusName.DECLINED) {
            throw new UserRejectedInvitationException(userId, classId);
        }
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
                        "Invitation", Map.of("classId", classId, "userId", userId)));
    }
}
