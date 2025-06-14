package org.example.userservice.security;

import lombok.RequiredArgsConstructor;
import org.example.userservice.entity.Role;
import org.example.userservice.entity.RoleName;
import org.example.userservice.entity.User;
import org.example.userservice.exception.OAuth2AuthenticationProcessingException;
import org.example.userservice.repository.RoleRepository;
import org.example.userservice.security.userinfo.OAuth2UserInfo;
import org.example.userservice.security.userinfo.OAuth2UserInfoFactory;
import org.example.userservice.service.UserService;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import java.util.Set;

@Service
@RequiredArgsConstructor
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    private final UserService userService;
    private final RoleRepository roleRepository;

    @Override
    @Transactional
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        OAuth2User oAuth2User = super.loadUser(userRequest);

        String provider = userRequest.getClientRegistration().getRegistrationId();
        OAuth2UserInfo oAuth2UserInfo = OAuth2UserInfoFactory.getOAuth2UserInfo(provider, oAuth2User.getAttributes());

        if (!StringUtils.hasText(oAuth2UserInfo.getEmail())) {
            throw new OAuth2AuthenticationProcessingException("Email not found from OAuth2 provider");
        }

        User user = userService
                .findByProviderAndProviderId(oAuth2UserInfo.getProvider(), oAuth2UserInfo.getId())
                .orElseGet(() -> registerNewUser(oAuth2UserInfo));

        return CustomUserDetails.create(user, oAuth2User.getAttributes());
    }

    private User registerNewUser(OAuth2UserInfo oAuth2UserInfo) {
        User newUser = new User();
        newUser.setProvider(oAuth2UserInfo.getProvider());
        newUser.setProviderId(oAuth2UserInfo.getId());
        newUser.setEmail(oAuth2UserInfo.getEmail());

        Role userRole = roleRepository.findByName(RoleName.USER)
                .orElseThrow(() -> new RuntimeException("Error: Default role USER not found in database."));
        newUser.setRoles(Set.of(userRole));

        return userService.save(newUser);
    }
}