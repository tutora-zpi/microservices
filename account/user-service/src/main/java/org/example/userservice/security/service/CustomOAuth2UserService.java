package org.example.userservice.security.service;

import lombok.RequiredArgsConstructor;
import org.example.userservice.entity.Role;
import org.example.userservice.entity.RoleName;
import org.example.userservice.entity.User;
import org.example.userservice.exception.OAuth2AuthenticationProcessingException;
import org.example.userservice.repository.RoleRepository;
import org.example.userservice.security.CustomUserDetails;
import org.example.userservice.security.userinfo.OAuth2UserInfo;
import org.example.userservice.security.userinfo.OAuth2UserInfoFactory;
import org.example.userservice.service.UserService;
import org.springframework.http.HttpEntity;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpMethod;
import org.springframework.http.ResponseEntity;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;
import org.springframework.web.client.RestTemplate;

import java.util.List;
import java.util.Map;
import java.util.Set;

@Service
@RequiredArgsConstructor
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    private final UserService userService;
    private final RoleRepository roleRepository;
    private final RestTemplate restTemplate;

    private static final String GITHUB_EMAILS_URL = "https://api.github.com/user/emails";

    @Override
    @Transactional
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        OAuth2User oAuth2User = super.loadUser(userRequest);

        String provider = userRequest.getClientRegistration().getRegistrationId();

        Map<String, Object> attributes = oAuth2User.getAttributes();
        OAuth2UserInfo oAuth2UserInfo = OAuth2UserInfoFactory.getOAuth2UserInfo(provider, attributes);

        String email = oAuth2UserInfo.getEmail();

        if (!StringUtils.hasText(email) && "github".equalsIgnoreCase(provider)) { //proxy email
            String githubId = oAuth2UserInfo.getId();
            email = githubId + "@users.noreply.github.com";
        }

        if (!StringUtils.hasText(email)) {
            throw new OAuth2AuthenticationProcessingException("Email not found from OAuth2 provider");
        }

        final String finalEmail = email;
        String providerId = oAuth2UserInfo.getId();

        User user = userService
                .findByProviderAndProviderId(provider, providerId)
                .orElseGet(() -> registerNewUser(provider, providerId, finalEmail));

        return CustomUserDetails.create(user, oAuth2User);
    }

    private String fetchPrimaryGitHubEmail(String accessToken) {
        HttpHeaders headers = new HttpHeaders();
        headers.setBearerAuth(accessToken);
        HttpEntity<String> entity = new HttpEntity<>("", headers);

        ResponseEntity<List<Map<String, Object>>> response = restTemplate.exchange(
                GITHUB_EMAILS_URL,
                HttpMethod.GET,
                entity,
                new org.springframework.core.ParameterizedTypeReference<>() {}
        );

        List<Map<String, Object>> emails = response.getBody();
        if (emails == null || emails.isEmpty()) {
            throw new OAuth2AuthenticationProcessingException("Email list not found from GitHub");
        }

        return emails.stream()
                .filter(emailMap -> (Boolean) emailMap.get("primary") && (Boolean) emailMap.get("verified"))
                .findFirst()
                .map(emailMap -> (String) emailMap.get("email"))
                .orElseThrow(() -> new OAuth2AuthenticationProcessingException("No primary verified email found from GitHub"));
    }

    private User registerNewUser(String provider, String id, String email) {
        User newUser = new User();
        newUser.setProvider(provider);
        newUser.setProviderId(id);
        newUser.setEmail(email);

        Role userRole = roleRepository.findByName(RoleName.USER)
                .orElseThrow(() -> new RuntimeException("Error: Default role USER not found in database."));
        newUser.setRoles(Set.of(userRole));

        return userService.save(newUser);
    }
}