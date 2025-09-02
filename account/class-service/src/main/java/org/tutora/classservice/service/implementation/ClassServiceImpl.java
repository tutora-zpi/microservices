package org.tutora.classservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.tutora.classservice.entity.Classroom;
import org.tutora.classservice.entity.Role;
import org.tutora.classservice.entity.RoleName;
import org.tutora.classservice.entity.UserClass;
import org.tutora.classservice.exception.ResourceNotFoundException;
import org.tutora.classservice.repository.ClassRepository;
import org.tutora.classservice.repository.RoleRepository;
import org.tutora.classservice.repository.UserClassRepository;
import org.tutora.classservice.service.contract.ClassService;

import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class ClassServiceImpl implements ClassService {

    private final ClassRepository classRepository;
    private final UserClassRepository userClassRepository;
    private final RoleRepository roleRepository;

    @Override
    public Classroom getClassById(UUID id) {
        return classRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Class", "id", id.toString()));
    }

    @Override
    public List<Classroom> getUserClasses(UUID userId) {
        return userClassRepository.findUserClassByUserId(userId).stream()
                .map(UserClass::getClassroom)
                .toList();
    }

    @Override
    @Transactional
    public Classroom createClass(UUID userId, Classroom newClassroom) {
        UserClass userClass = new UserClass();
        Role role = roleRepository.findByName(RoleName.HOST)
                        .orElseThrow(() -> new ResourceNotFoundException("Role", "name", RoleName.HOST.toString()));

        userClass.setClassroom(newClassroom);
        userClass.setUserId(userId);
        userClass.setRole(role);

        Classroom savedClassroom = classRepository.save(newClassroom);
        userClassRepository.save(userClass);

        return savedClassroom;
    }
}
