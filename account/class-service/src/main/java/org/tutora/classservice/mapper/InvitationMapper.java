package org.tutora.classservice.mapper;

import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.tutora.classservice.dto.InvitationDto;
import org.tutora.classservice.entity.Invitation;

@Mapper(componentModel = "spring")
public interface InvitationMapper {

    @Mapping(target = "status", source = "status", qualifiedByName = "mapStatusesToString")
    InvitationDto toDto(Invitation invitation);

    @Named("mapStatusesToString")
    default String mapStatusesToString(final Invitation invitation) {
        return invitation.getStatus().toString();
    }
}
