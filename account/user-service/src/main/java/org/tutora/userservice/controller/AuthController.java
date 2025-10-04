package org.tutora.userservice.controller;

import lombok.RequiredArgsConstructor;
import org.tutora.userservice.dto.UserDto;
import org.tutora.userservice.mapper.UserMapper;
import org.tutora.userservice.security.CustomUserDetails;
import org.tutora.userservice.service.contract.UserService;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/auth")
@RequiredArgsConstructor
public class AuthController {

    private final UserService userService;
    private final UserMapper userMapper;

    /**
     * Endpoint zwracający dane o aktualnie zalogowanym użytkowniku.
     * Dostęp do niego jest chroniony przez JwtAuthenticationFilter.
     * Wymaga przesłania ważnego tokenu JWT w nagłówku "Authorization: Bearer <token>".
     *
     * @param userDetails Obiekt principal (CustomUserDetails) wstrzyknięty przez Spring Security
     * po pomyślnej walidacji tokenu JWT.
     * @return Dane zalogowanego użytkownika w formacie DTO.
     */
    @GetMapping("/me")
    public ResponseEntity<UserDto> getCurrentUser(@AuthenticationPrincipal CustomUserDetails userDetails) {
        var user = userService.findById(userDetails.getId());

        return ResponseEntity.ok(userMapper.toDto(user));
    }
}