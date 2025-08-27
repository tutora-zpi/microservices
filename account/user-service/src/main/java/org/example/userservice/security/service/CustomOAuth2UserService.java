package org.example.userservice.security.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.example.userservice.entity.User;
import org.example.userservice.exception.OAuth2AuthenticationProcessingException;
import org.example.userservice.security.CustomUserDetails;
import org.example.userservice.security.userinfo.OAuth2UserInfo;
import org.example.userservice.security.userinfo.OAuth2UserInfoFactory;
import org.example.userservice.service.contract.UserService;
import org.springframework.security.oauth2.client.userinfo.DefaultOAuth2UserService;
import org.springframework.security.oauth2.client.userinfo.OAuth2UserRequest;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.user.OAuth2User;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

import java.util.Map;

@Service
@RequiredArgsConstructor
@Slf4j
public class CustomOAuth2UserService extends DefaultOAuth2UserService {

    private final UserService userService;

    @Override
    @Transactional
    public OAuth2User loadUser(OAuth2UserRequest userRequest) throws OAuth2AuthenticationException {
        log.info(userRequest.toString());

        OAuth2User oAuth2User = super.loadUser(userRequest);

        log.info(oAuth2User.toString());

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
                .orElseGet(() -> userService.registerUser(oAuth2UserInfo));

        return CustomUserDetails.create(user, oAuth2User);
    }
}