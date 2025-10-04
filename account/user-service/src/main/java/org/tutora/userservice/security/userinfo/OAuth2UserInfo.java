package org.tutora.userservice.security.userinfo;

public interface OAuth2UserInfo {
    String getId();
    String getEmail();
    String getProvider();
    String getName();
    String getSurname();
    String getImageUrl();
}
