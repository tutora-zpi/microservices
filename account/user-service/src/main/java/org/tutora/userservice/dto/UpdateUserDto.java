package org.tutora.userservice.dto;

import lombok.Data;

@Data
public class UpdateUserDto {
    private String email;
    private String name;
    private String surname;
}
