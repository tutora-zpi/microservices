package org.tutora.userservice.security.service;

import lombok.RequiredArgsConstructor;
import org.tutora.userservice.entity.User;
import org.tutora.userservice.exception.OAuth2AuthenticationProcessingException;
import org.tutora.userservice.security.CustomUserDetails;
import org.tutora.userservice.security.userinfo.OAuth2UserInfo;
import org.tutora.userservice.security.userinfo.OAuth2UserInfoFactory;
import org.tutora.userservice.service.contract.UserService;
import org.springframework.security.oauth2.client.oidc.userinfo.OidcUserRequest;
import org.springframework.security.oauth2.client.oidc.userinfo.OidcUserService;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.core.oidc.user.OidcUser;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;

@Service
@RequiredArgsConstructor
public class CustomOidcUserService extends OidcUserService {

    private final UserService userService;

    @Override
    @Transactional
    public OidcUser loadUser(OidcUserRequest userRequest) throws OAuth2AuthenticationException {
        OidcUser oidcUser = super.loadUser(userRequest);

        String provider = userRequest.getClientRegistration().getRegistrationId();
        OAuth2UserInfo oAuth2UserInfo = OAuth2UserInfoFactory.getOAuth2UserInfo(provider, oidcUser.getAttributes());

        if (!StringUtils.hasText(oAuth2UserInfo.getEmail())) {
            throw new OAuth2AuthenticationProcessingException("Email not found from OIDC provider");
        }

        User user = userService
                .findByProviderAndProviderId(oAuth2UserInfo.getProvider(), oAuth2UserInfo.getId())
                .orElseGet(() -> userService.registerUser(oAuth2UserInfo));

        return CustomUserDetails.create(user, oidcUser);
    }
}