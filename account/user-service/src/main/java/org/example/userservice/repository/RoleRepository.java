package org.example.userservice.repository;

import org.example.userservice.entity.Role;
import org.example.userservice.entity.RoleName;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;

public interface RoleRepository extends JpaRepository<Role, Long> {
    Optional<Role> findByName(RoleName roleName);
}
