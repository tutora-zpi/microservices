package org.tutora.classservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.tutora.classservice.entity.Member;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

public interface MemberRepository extends JpaRepository<Member, Long> {
    List<Member> findUserClassByUserId(UUID userId);
    Optional<Member> findUserClassByUserIdAndClassroomId(UUID userId, UUID classId);
    boolean existsByClassroomIdAndUserId(UUID classId, UUID userId);
}
