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
import org.example.userservice.service.contract.AvatarService;
import org.example.userservice.service.contract.UserService;
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
    private final AvatarService avatarService;
    private final RoleRepository roleRepository;

    @Override
    @Transactional
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        OAuth2User oAuth2User = super.loadUser(userRequest);

        String provider = userRequest.getClientRegistration().getRegistrationId();

        Map<String, Object> attributes = oAuth2User.getAttributes();
        OAuth2UserInfo oAuth2UserInfo = OAuth2UserInfoFactory.getOAuth2UserInfo(provider, attributes);

        String email = oAuth2UserInfo.getEmail();

        if (!StringUtils.hasText(email)) {
            throw new OAuth2AuthenticationProcessingException("Email not found from OAuth2 provider");
        }

        String providerId = oAuth2UserInfo.getId();

        User user = userService
                .findByProviderAndProviderId(provider, providerId)
                .orElseGet(() -> registerNewUser(oAuth2UserInfo));

        return CustomUserDetails.create(user, oAuth2User);
    }

    private User registerNewUser(OAuth2UserInfo oAuth2UserInfo) {
        User newUser = new User();
        newUser.setProvider(oAuth2UserInfo.getProvider());
        newUser.setProviderId(oAuth2UserInfo.getId());
        newUser.setEmail(oAuth2UserInfo.getEmail());
        newUser.setName(oAuth2UserInfo.getName());
        newUser.setSurname(oAuth2UserInfo.getSurname());

        Role userRole = roleRepository.findByName(RoleName.USER)
                .orElseThrow(() -> new RuntimeException("Error: Default role USER not found in database."));
        newUser.setRoles(Set.of(userRole));

        String pictureUrl = oAuth2UserInfo.getImageUrl();
        String avatarKey = avatarService.saveAvatarFromUrl(newUser.getId(), pictureUrl);

        newUser.setAvatarKey(avatarKey);

        return userService.save(newUser);
    }
}