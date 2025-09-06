package org.tutora.classservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.tutora.classservice.entity.UserClass;

import java.util.List;
import java.util.UUID;

public interface UserClassRepository extends JpaRepository<UserClass, Long> {
    List<UserClass> findUserClassByUserId(UUID userId);
    boolean existsByClassroomIdAndUserId(UUID classId, UUID userId);
}
