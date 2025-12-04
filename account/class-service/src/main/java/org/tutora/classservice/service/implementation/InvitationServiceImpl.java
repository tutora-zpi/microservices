package org.tutora.classservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.tutora.classservice.client.NotificationPublisher;
import org.tutora.classservice.dto.UserDto;
import org.tutora.classservice.event.ClassInvitationAcceptedEvent;
import org.tutora.classservice.event.ClassInvitationCreatedEvent;
import org.tutora.classservice.entity.*;
import org.tutora.classservice.exception.*;
import org.tutora.classservice.mapper.ClassMapper;
import org.tutora.classservice.repository.InvitationRepository;
import org.tutora.classservice.repository.InvitationStatusRepository;
import org.tutora.classservice.repository.MemberRepository;
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
    private final MemberRepository memberRepository;
    private final ClassMapper classMapper;

    @Override
    public Invitation inviteUser(UserDto sender, UUID classId, UserDto receiver) {
        Classroom classroom = classService.getClassById(classId);

        UUID senderId = UUID.fromString(sender.id());
        if (!hasAuthority(senderId, classId, RoleName.HOST)) {
            throw new UnauthorizedActionException("classroom", classId, "send invitation to classroom");
        }

        UUID receiverId = UUID.fromString(receiver.id());
        if (memberRepository.existsByClassroomIdAndUserId(classId, receiverId)) {
            throw new UserAlreadyInClassException(receiverId, classId);
        }

        Optional<Invitation> existingInvitation = invitationRepository.findByClassroomIdAndUserId(classId, receiverId);
        if(existingInvitation.isPresent()) {
            Invitation inv = existingInvitation.get();
            switch (inv.getStatus().getStatusName()) {
                case DECLINED -> throw new UserRejectedInvitationException(receiverId, classId);
                case INVITED -> throw new UserAlreadyInvitedException(receiverId, classId);
            }
        }

        Invitation inv = saveInvitation(classroom, receiverId);

        notificationPublisher.sendClassInvitation(new ClassInvitationCreatedEvent(
                classroom.getName(),
                receiver,
                sender
        ));

        return inv;
    }

    @Transactional
    @Override
    public void cancelInvitation(UUID classId, UUID userId) {
        Invitation inv = getInvitation(classId, userId);

        validateInvitationStatus(inv.getStatus(), classId, userId);

        invitationRepository.deleteByClassroomIdAndUserId(classId, userId);
    }

    @Override
    public void joinClass(UUID classId, UUID userId) {
        Invitation inv = getInvitation(classId, userId);

        validateInvitationStatus(inv.getStatus(), classId, userId);

        inv.setStatus(getInvitationStatus(InvitationStatusName.ACCEPTED));
        invitationRepository.save(inv);

        classService.addUserToClass(classId, userId, RoleName.GUEST);

        Classroom classroom = classService.getClassById(classId);

        notificationPublisher.sendClassInvitationAccepted(new ClassInvitationAcceptedEvent(
                classId,
                classroom.getName(),
                classService.getClassHost(classId).getUserId(),
                userId
        ));
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

    private boolean hasAuthority(UUID userId, UUID classId, RoleName role) {
        Member member = memberRepository.findMemberByUserIdAndClassroomId(userId, classId)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "Member", Map.of("userId", userId, "classId", classId)));

        return member.getRole().getName() == role;
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

    private Invitation saveInvitation(Classroom classroom, UUID userId) {
        InvitationStatus status = getInvitationStatus(InvitationStatusName.INVITED);

        Invitation inv = Invitation.builder()
                .userId(userId)
                .status(status)
                .createdAt(LocalDateTime.now())
                .build();
        classroom.addInvitation(inv);

        return invitationRepository.save(inv);
    }

    private Invitation getInvitation(UUID classId, UUID userId) {
        return invitationRepository.findByClassroomIdAndUserId(classId, userId)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "Invitation", Map.of("classId", classId, "userId", userId)));
    }
}
