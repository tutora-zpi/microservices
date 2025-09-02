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
        Classroom savedClassroom = classRepository.save(newClassroom);
        saveUserClass(userId, savedClassroom, RoleName.HOST);

        return savedClassroom;
    }

    @Override
    public void addUserToClass(UUID classId, UUID userId, RoleName role) {
        Classroom classroom = getClassById(classId);

        saveUserClass(userId, classroom, role);
    }

    private Role getRoleByName(RoleName name) {
        return roleRepository.findByName(name)
                .orElseThrow(() -> new ResourceNotFoundException("Role", "name", name.toString()));
    }

    private void saveUserClass(UUID userId, Classroom classroom, RoleName name) {
        Role role = getRoleByName(name);

        UserClass userClass = UserClass.builder()
                .classroom(classroom)
                .userId(userId)
                .role(role)
                .build();

        userClassRepository.save(userClass);
    }
}
