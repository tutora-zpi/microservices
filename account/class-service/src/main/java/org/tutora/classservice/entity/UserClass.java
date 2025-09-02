package org.tutora.classservice.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.annotations.OnDelete;
import org.hibernate.annotations.OnDeleteAction;

import java.util.UUID;

@Entity
@Table(
        name = "users_classes"
)
@Getter
@Setter
@AllArgsConstructor
@NoArgsConstructor
public class UserClass {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @ManyToOne(cascade = CascadeType.MERGE)
    @JoinColumn(
            name = "class_id"
    )
    @OnDelete(action = OnDeleteAction.CASCADE)
    private Classroom classroom;

    @Column(nullable = false)
    private UUID userId;

    @ManyToOne(cascade = CascadeType.MERGE)
    @JoinColumn(
            name = "user_role"
    )
    private Role role;
}
