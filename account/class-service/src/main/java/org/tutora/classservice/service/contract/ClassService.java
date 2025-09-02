package org.tutora.classservice.service.contract;

import org.tutora.classservice.entity.Classroom;

import java.util.List;
import java.util.UUID;

public interface ClassService {
    Classroom getClassById(UUID id);
    List<Classroom> getUserClasses(UUID userId);
    Classroom createClass(UUID userId, Classroom newClassroom);
}
