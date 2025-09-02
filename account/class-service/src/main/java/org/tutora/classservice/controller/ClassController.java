package org.tutora.classservice.controller;


import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;
import org.springframework.web.bind.annotation.*;
import org.tutora.classservice.dto.ClassCreateRequest;
import org.tutora.classservice.dto.ClassDto;
import org.tutora.classservice.entity.Classroom;
import org.tutora.classservice.mapper.ClassMapper;
import org.tutora.classservice.service.contract.AuthService;
import org.tutora.classservice.service.contract.ClassService;

import java.util.List;
import java.util.UUID;

@RestController
@RequestMapping("/classes")
@RequiredArgsConstructor
@Slf4j
public class ClassController {

    private final AuthService authService;
    private final ClassService classService;
    private final ClassMapper classMapper;

    @PostMapping
    public ResponseEntity<ClassDto> createClass(@RequestBody ClassCreateRequest classCreateRequest) {
        log.info("Request to create class: {}", classCreateRequest);

        Classroom classroom = classService.createClass(
                classCreateRequest.hostId(),
                classMapper.toEntity(classCreateRequest)
        );

        return ResponseEntity
                .status(HttpStatus.CREATED)
                .body(classMapper.toDto(classroom));
    }

    @GetMapping
    public ResponseEntity<List<ClassDto>> getUserClasses(@AuthenticationPrincipal Jwt principal) {
        log.info("Request to get user classes for user with principal: {}", principal);
        UUID userId = UUID.fromString(authService.getUserId(principal));

        List<ClassDto> userClasses = classService.getUserClasses(userId).stream()
                .map(classMapper::toDto)
                .toList();

        return ResponseEntity
                .ok(userClasses);
    }

    @GetMapping("/{id}")
    public ResponseEntity<ClassDto> getClass(@PathVariable UUID id) {
        log.info("Request to get class with id: {}", id);

        Classroom classroom = classService.getClassById(id);

        return ResponseEntity
                .ok(classMapper.toDto(classroom));
    }
}
