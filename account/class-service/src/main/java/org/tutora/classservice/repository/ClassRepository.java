package org.tutora.classservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import org.tutora.classservice.entity.Classroom;

import java.util.List;
import java.util.Optional;
import java.util.UUID;

@Repository
public interface ClassRepository extends JpaRepository<Classroom, UUID> {
    Optional<Classroom> findClassroomById(UUID id);
}