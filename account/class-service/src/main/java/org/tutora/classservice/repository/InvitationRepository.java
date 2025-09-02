package org.tutora.classservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import org.tutora.classservice.entity.Invitation;
import org.tutora.classservice.entity.InvitationStatus;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface InvitationRepository extends JpaRepository<Invitation, Integer> {
    Optional<Invitation> findByClassroomIdAndUserId(UUID classId, UUID userId);
    List<Invitation> findAllByUserIdAndStatus(UUID userId, InvitationStatus status);
    List<Invitation> findAllByClassroomIdAndStatus(UUID classId, InvitationStatus status);
    void deleteByClassroomIdAndUserId(UUID classId, UUID userId);
}
