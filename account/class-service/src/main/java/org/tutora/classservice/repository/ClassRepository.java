package org.tutora.classservice.repository;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;
import org.tutora.classservice.entity.Class;

import java.util.UUID;

@Repository
public interface ClassRepository extends JpaRepository<Class, UUID> {
}