package org.example.userservice.mapper;

import org.example.userservice.dto.UpdateUserDto;
import org.example.userservice.dto.UserDto;
import org.example.userservice.entity.Role;
import org.example.userservice.entity.User;
import org.example.userservice.service.contract.AvatarService;
import org.mapstruct.*;
import org.springframework.beans.factory.annotation.Autowired;

import java.util.Set;
import java.util.stream.Collectors;

@Mapper(componentModel = "spring")
public abstract class UserMapper {

    @Autowired
    protected AvatarService avatarService;

    @Mapping(target = "roles", source = "roles", qualifiedByName = "mapRolesToNames")
    @Mapping(target = "avatarUrl", source = "avatarKey", qualifiedByName = "mapAvatarKeyToUrl")
    public abstract UserDto toDto(User user);

    @Named("mapRolesToNames")
    protected Set<String> mapRolesToNames(Set<Role> roles) {
        if (roles == null) {
            return Set.of();
        }
        return roles.stream()
                .map(role -> role.getName().name())
                .collect(Collectors.toSet());
    }

    @Named("mapAvatarKeyToUrl")
    protected String mapAvatarKeyToUrl(String avatarKey) {
        if (avatarKey == null) return null;
        return avatarService.getAvatarUrl(avatarKey);
    }
}
