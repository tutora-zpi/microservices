package org.tutora.classservice.service.contract;

import org.tutora.classservice.entity.Invitation;

import java.util.List;
import java.util.UUID;

public interface InvitationService {
    Invitation inviteUser(String senderFullName, UUID classId, UUID receiverId);
    void cancelInvitation(UUID classId, UUID receiverId);
    void joinClass(UUID classId, UUID receiverId);
    void declineInvitation(UUID classId, UUID receiverId);
    List<Invitation> getInvitationsForUser(UUID userId);
    List<Invitation> getInvitationsForClass(UUID classId);
}
