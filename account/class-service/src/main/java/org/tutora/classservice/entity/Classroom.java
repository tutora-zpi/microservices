package org.tutora.classservice.entity;

import jakarta.persistence.*;
import lombok.*;

import java.time.LocalDateTime;
import java.util.ArrayList;
import java.util.List;
import java.util.UUID;

@Entity
@Table(
        name = "classes"
)
@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
@Builder
public class Classroom {

    @Id
    @GeneratedValue(strategy = GenerationType.UUID)
    private UUID id;

    @Column(nullable = false)
    private String name;

    @Column
    private LocalDateTime createdAt;

    @OneToMany(mappedBy = "classroom", cascade = CascadeType.ALL, orphanRemoval = true)
    @Builder.Default
    private List<UserClass> userClasses = new ArrayList<>();

    public void addUserClass(UserClass userClass) {
        userClasses.add(userClass);
        userClass.setClassroom(this);
    }
}
