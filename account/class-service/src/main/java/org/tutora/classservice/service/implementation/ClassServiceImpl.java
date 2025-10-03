package org.tutora.classservice.service.implementation;

import lombok.RequiredArgsConstructor;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.tutora.classservice.entity.Classroom;
import org.tutora.classservice.entity.Role;
import org.tutora.classservice.entity.RoleName;
import org.tutora.classservice.entity.Member;
import org.tutora.classservice.exception.ResourceNotFoundException;
import org.tutora.classservice.exception.UnauthorizedActionException;
import org.tutora.classservice.repository.ClassRepository;
import org.tutora.classservice.repository.RoleRepository;
import org.tutora.classservice.repository.MemberRepository;
import org.tutora.classservice.service.contract.ClassService;

import java.util.List;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class ClassServiceImpl implements ClassService {

    private final ClassRepository classRepository;
    private final MemberRepository memberRepository;
    private final RoleRepository roleRepository;

    @Override
    public Classroom getClassById(UUID id) {
        return classRepository.findById(id)
                .orElseThrow(() -> new ResourceNotFoundException("Class", "id", id.toString()));
    }

    @Override
    public List<Classroom> getUserClasses(UUID userId) {
        return memberRepository.findMemberByUserId(userId).stream()
                .map(Member::getClassroom)
                .toList();
    }

    @Override
    @Transactional
    public Classroom createClass(UUID userId, Classroom newClassroom) {
        newClassroom.setId(null);
        Member member = createUserClass(userId, RoleName.HOST);
        newClassroom.addUserClass(member);

        return classRepository.save(newClassroom);
    }

    @Override
    public void addUserToClass(UUID classId, UUID userId, RoleName role) {
        Classroom classroom = getClassById(classId);
        classroom.addUserClass(createUserClass(userId, role));

        classRepository.save(classroom);
    }

    @Override
    public void deleteClass(UUID classId, UUID userID) {
        Member member = memberRepository.findMemberByUserIdAndClassroomId(userID, classId)
                .orElseThrow(() -> new ResourceNotFoundException("Member", "id", classId.toString()));

        if (member.getRole().getName() == RoleName.HOST) {
            classRepository.deleteById(classId);
        } else {
            throw new UnauthorizedActionException("class", classId, "delete");
        }
    }

    private Role getRoleByName(RoleName name) {
        return roleRepository.findByName(name)
                .orElseThrow(() -> new ResourceNotFoundException("Role", "name", name.toString()));
    }

    private Member createUserClass(UUID userId, RoleName name) {
        Role role = getRoleByName(name);

        return Member.builder()
                .userId(userId)
                .role(role)
                .build();
    }
}
