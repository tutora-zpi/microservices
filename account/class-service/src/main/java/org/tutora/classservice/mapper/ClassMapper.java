package org.tutora.classservice.mapper;

import org.mapstruct.*;
import org.tutora.classservice.dto.ClassCreateRequest;
import org.tutora.classservice.dto.ClassDto;
import org.tutora.classservice.dto.UserRoleDto;
import org.tutora.classservice.entity.Classroom;
import org.tutora.classservice.entity.Member;

import java.time.LocalDateTime;
import java.util.List;
import java.util.UUID;

@Mapper(componentModel = "spring")
public interface ClassMapper {

    @Mapping(target = "id", expression = "java(generateUuid())")
    @Mapping(target = "createdAt", expression = "java(currentTime())")
    Classroom toEntity(ClassCreateRequest dto);

    default UUID generateUuid() {
        return UUID.randomUUID();
    }

    default LocalDateTime currentTime() {
        return LocalDateTime.now();
    }

    @Mapping(target = "members", source = "members", qualifiedByName = "mapMembers")
    ClassDto toDto(Classroom classroom);

    @Named("mapMembers")
    default List<UserRoleDto> mapMembers(List<Member> members) {
        if (members == null)
            return null;

        return members.stream()
                .map(uc -> new UserRoleDto(uc.getUserId(), uc.getRole().getName().name()))
                .toList();
    }
}
