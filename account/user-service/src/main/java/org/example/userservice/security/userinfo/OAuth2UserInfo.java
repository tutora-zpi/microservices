package org.example.userservice.security.userinfo;

public interface OAuth2UserInfo {
    String getId();
    String getEmail();
    String getProvider();
}
