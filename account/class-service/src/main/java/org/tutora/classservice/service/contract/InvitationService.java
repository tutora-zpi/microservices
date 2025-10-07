package org.tutora.classservice.service.contract;

import org.tutora.classservice.dto.UserDto;
import org.tutora.classservice.entity.Invitation;

import java.util.List;
import java.util.UUID;

public interface InvitationService {
    Invitation inviteUser(UserDto sender, UUID classId, UserDto receiver);
    void cancelInvitation(UUID classId, UUID receiverId);
    void joinClass(UUID classId, UUID receiverId);
    void declineInvitation(UUID classId, UUID receiverId);
    List<Invitation> getInvitationsForUser(UUID userId);
    List<Invitation> getInvitationsForClass(UUID classId);
}
