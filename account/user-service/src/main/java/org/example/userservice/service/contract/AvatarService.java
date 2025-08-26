package org.example.userservice.service.contract;

import org.example.userservice.entity.User;

public interface AvatarService {

    /**
     * Generuje pre-signed URL do wrzucenia avatara użytkownika
     *
     * @param key       - klucz obiektu w storage (np. avatars/123.png)
     * @param contentType  - typ MIME pliku (np. image/png, image/jpeg)
     * @return pre-signed URL do wrzucenia pliku
     */
    String generateUploadUrl(String key, String contentType);

    /// DEFAULT AVATAR JEŚLI BRAK AVATARA?
    /**
     * Zwraca URL do pobrania avatara (np. przez CloudFront).
     * Zwraca null jeśli brak avatara.
     *
     * @param key - klucz obiektu w storage (np. avatars/123.png)
     * @return URL do pobrania
     */
    String getAvatarUrl(String key);

    /**
     * Usuwa avatar ze storage
     *
     * @param key - klucz obiektu w storage (np. avatars/123.png)
     */
    void deleteAvatar(String key);
}
