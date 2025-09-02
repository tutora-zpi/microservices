package org.tutora.classservice.service.contract;

import org.tutora.classservice.entity.Invitation;

import java.util.List;
import java.util.UUID;

public interface InvitationService {
    void inviteUser(UUID classId, UUID userId);
    void cancelInvitation(UUID classId, UUID userId);
    void joinClass(UUID classId, UUID userId);
    void declineInvitation(UUID classId, UUID userId);
    List<Invitation> getInvitationsForUser(UUID userId);
    List<Invitation> getInvitationsForClass(UUID classId);
}
