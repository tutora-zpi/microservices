package org.example.userservice.service;

public interface AvatarService {

    /**
     * Generuje pre-signed URL do wrzucenia avatara użytkownika
     *
     * @param userId       - ID użytkownika
     * @param contentType  - typ MIME pliku (np. image/png, image/jpeg)
     * @return pre-signed URL do wrzucenia pliku
     */
    String generateUploadUrl(Long userId, String contentType);

    /**
     * Zwraca URL do pobrania avatara (np. przez CloudFront).
     * Może zwrócić null lub URL defaultowego avatara jeśli brak avatara.
     *
     * @param avatarKey - klucz obiektu w storage (np. avatars/123.png)
     * @return URL do pobrania
     */
    String getAvatarUrl(String avatarKey);
}
