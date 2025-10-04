package org.tutora.userservice.repository;

import org.tutora.userservice.entity.Role;
import org.tutora.userservice.entity.RoleName;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;

public interface RoleRepository extends JpaRepository<Role, Long> {
    Optional<Role> findByName(RoleName roleName);
}
