package org.tutora.userservice.repository;

import org.tutora.userservice.entity.User;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;

import java.util.Optional;
import java.util.UUID;

public interface UserRepository extends JpaRepository<User, UUID> {
    Optional<User> findByProviderAndProviderId(String provider, String providerId);
    Page<User> findByNameStartingWithIgnoreCaseOrSurnameStartingWithIgnoreCase(
            String firstName,
            String lastName,
            Pageable pageable
    );
}
