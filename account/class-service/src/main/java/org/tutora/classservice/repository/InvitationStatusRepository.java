package org.tutora.classservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.tutora.classservice.entity.InvitationStatus;
import org.tutora.classservice.entity.InvitationStatusName;

import java.util.Optional;

public interface InvitationStatusRepository extends JpaRepository<InvitationStatus, Integer> {
    Optional<InvitationStatus> findByStatusName(InvitationStatusName status);
}
