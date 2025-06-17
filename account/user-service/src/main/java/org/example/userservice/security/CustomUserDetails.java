package org.example.userservice.security;

import lombok.Getter;
import org.example.userservice.entity.User;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.SimpleGrantedAuthority;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.oauth2.core.oidc.OidcIdToken;
import org.springframework.security.oauth2.core.oidc.OidcUserInfo;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.security.oauth2.core.user.OAuth2User;

import java.util.Collection;
import java.util.Map;
import java.util.UUID;
import java.util.stream.Collectors;

@Getter
public class CustomUserDetails implements OidcUser, UserDetails {

    // Dane z naszej bazy danych
    private final UUID id;
    private final String email;
    private final Collection<? extends GrantedAuthority> authorities;

    // Dane od dostawcy OAuth2/OIDC
    private final Map<String, Object> attributes;
    private final OidcIdToken idToken;
    private final OidcUserInfo userInfo;

    /**
     * Prywatny konstruktor, aby wymusić tworzenie obiektu przez metody fabryczne.
     */
    private CustomUserDetails(User user, Map<String, Object> attributes, OidcIdToken idToken, OidcUserInfo userInfo) {
        this.id = user.getId();
        this.email = user.getEmail();
        this.authorities = user.getRoles().stream()
                .map(role -> new SimpleGrantedAuthority("ROLE_" + role.getName().name()))
                .collect(Collectors.toList());
        this.attributes = attributes;
        this.idToken = idToken;
        this.userInfo = userInfo;
    }

    /**
     * Metoda fabryczna dla użytkowników OIDC (np. Google).
     * @param user Nasza encja użytkownika z bazy danych.
     * @param oidcUser Oryginalny obiekt użytkownika od dostawcy OIDC.
     * @return Nowa instancja CustomUserDetails.
     */
    public static CustomUserDetails create(User user, OidcUser oidcUser) {
        return new CustomUserDetails(
                user,
                oidcUser.getAttributes(),
                oidcUser.getIdToken(),
                oidcUser.getUserInfo()
        );
    }

    /**
     * Metoda fabryczna dla "zwykłych" użytkowników OAuth2 (np. GitHub).
     * @param user Nasza encja użytkownika z bazy danych.
     * @param oauth2User Oryginalny obiekt użytkownika od dostawcy OAuth2.
     * @return Nowa instancja CustomUserDetails.
     */
    public static CustomUserDetails create(User user, OAuth2User oauth2User) {
        return new CustomUserDetails(
                user,
                oauth2User.getAttributes(),
                null,
                null
        );
    }

    /**
     * NOWA METODA FABRYCZNA - Dla JwtAuthenticationFilter
     * Tworzy obiekt Principal na podstawie danych z bazy, bez kontekstu OAuth2.
     * @param user Nasza encja użytkownika.
     * @return Nowa instancja CustomUserDetails.
     */
    public static CustomUserDetails create(User user) {
        return new CustomUserDetails(
                user,
                Map.of(), // Atrybuty są puste, bo nie pochodzą z OAuth2
                null,
                null
        );
    }

    @Override
    public String getUsername() { return email; }

    @Override
    public String getPassword() { return null; }

    @Override
    public boolean isAccountNonExpired() { return true; }

    @Override
    public boolean isAccountNonLocked() { return true; }

    @Override
    public boolean isCredentialsNonExpired() { return true; }

    @Override
    public boolean isEnabled() { return true; }

    @Override
    public Map<String, Object> getAttributes() { return this.attributes; }

    @Override
    public Collection<? extends GrantedAuthority> getAuthorities() { return this.authorities; }

    @Override
    public String getName() { return String.valueOf(this.id); }

    @Override
    public Map<String, Object> getClaims() { return this.attributes; }

    @Override
    public OidcUserInfo getUserInfo() {
        return this.userInfo;
    }

    @Override
    public OidcIdToken getIdToken() {
        return this.idToken;
    }
}