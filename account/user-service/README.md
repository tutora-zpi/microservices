# Serwis Autoryzacji (Auth-Service)

Centralny mikroserwis odpowiedzialny za uwierzytelnianie użytkowników w ekosystemie. Obsługuje logowanie przez dostawców zewnętrznych (OAuth2) i wydaje tokeny JWT podpisane asymetrycznie (RSA), które mogą być weryfikowane przez inne serwisy bez potrzeby odpytywania serwisu autoryzacji.

---

## Funkcjonalności

* **Logowanie jednokrotne (SSO)** z dostawcami:
    * Google (automatyczne pobieranie danych i zapisywanie w bazie)
* **Generowanie tokenów JWT** podpisanych asymetrycznie algorytmem **RS256**.
* **Trwałość użytkowników** w bazie danych PostgreSQL.
* **Zarządzanie danymi użytkowników**:
  * Pobieranie danych użytkownika po UUID
  * Aktualizacja danych użytkownika (`email`, `name`, `surname`)
  * Generowanie presigned URL do uploadu awatarów
* **Endpoint JWKS** (`/.well-known/jwks.json`) do dystrybucji klucza publicznego.
* **Zdecentralizowana walidacja tokenów** – każdy mikroserwis może weryfikować tokeny samodzielnie.

## Architektura

Serwis implementuje nowoczesny i bezpieczny wzorzec uwierzytelniania dla systemów rozproszonych. Po pomyślnej autoryzacji u dostawcy OAuth2, serwis tworzy lub aktualizuje użytkownika w swojej bazie danych, a następnie generuje token JWT podpisany kluczem prywatnym. Klucz publiczny jest udostępniany przez standardowy endpoint JWKS, co pozwala innym serwisom na bezstanową i szybką weryfikację tokenów.

Dodatkowo serwis zarządza pełnym cyklem życia użytkowników:
* Przechowuje ich dane w bazie PostgreSQL (`email`, `name`, `surname`, `role`, `avatarKey`).
* Umożliwia pobranie danych użytkownika po UUID.
* Umożliwia aktualizację podstawowych danych użytkownika (email, imię, nazwisko).
* Generuje presigned URL do uploadu awatarów i zapisuje referencję w bazie.
---

## Wymagania

Przed uruchomieniem upewnij się, że masz zainstalowane następujące narzędzia:
* Java 17 (lub nowsza)
* Apache Maven
* Docker i Docker Compose
* Skonfigurowany bucket w AWS S3
* `openssl` (do generowania kluczy RSA, zazwyczaj dostępne w systemach Linux/macOS lub przez Git Bash w Windows)

---

## Konfiguracja

### 1. Generowanie Kluczy RSA

Każde środowisko (deweloperskie i produkcyjne) powinno używać własnej, unikalnej pary kluczy.

Otwórz terminal w głównym katalogu projektu i wykonaj poniższe komendy. Stworzą one katalog `keys` z kluczami, który jest ignorowany przez Git.

```bash
mkdir -p src/main/resources/keys

# Generowanie klucza prywatnego
openssl genpkey -algorithm RSA -out src/main/resources/keys/private_key.pem -pkeyopt rsa_keygen_bits:2048

# Generowanie klucza publicznego
openssl rsa -pubout -in src/main/resources/keys/private_key.pem -out src/main/resources/keys/public_key.pem
```

### 2. Konfiguracja Dostawców OAuth2

Musisz założyć własne aplikacje OAuth2 w konsolach deweloperskich Google i GitHub, aby uzyskać `Client ID` i `Client Secret`.

* **Google:** [Google Cloud Console](https://console.cloud.google.com/apis/credentials)

Podczas konfiguracji aplikacji u dostawców, jako **Authorized redirect URI** podaj:
`http://localhost:8080/login/oauth2/code/google`

### 3. Utworzenie pliku `.env`


Stwórz plik `.env` na podstawie `.env.sample`

**Otwórz plik `.env` i uzupełnij go** swoimi danymi uzyskanymi w poprzednim kroku oraz danymi do bazy danych.

**WAŻNE:** Plik `.env` zawiera poufne dane i **nie powinien** być dodawany do systemu kontroli wersji Git. Jest on już uwzględniony w pliku `.gitignore`.

---

## Uruchomienie

Upewnij się, że jesteś w głównym katalogu projektu i wykonaj komendę:
```bash
docker-compose up --build
```


---

## Endpointy API

| Metoda  | Ścieżka | Opis | Dostęp | Body                                                  |
|:--------| :--- | :--- | :--- |:------------------------------------------------------|
| `GET`   | `/oauth2/authorization/{provider}` | Inicjuje proces logowania przez danego dostawcę (`google` lub `github`). | Publiczny | –                                                     |
| `GET`   | `/.well-known/jwks.json` | Zwraca klucz publiczny RSA w formacie JWKS do weryfikacji tokenów. | Publiczny | –                                                     |
| `POST`  | `/oauth2/token` | Generuje token JWT dla bota/serwisu w zamian za jego `client_id` i `client_secret`. | Publiczny | `x-www-form-urlencoded` (`grant_type=client_credentials`) |
| `GET`   | `/auth/me` | Zwraca dane o aktualnie zalogowanym użytkowniku. | Chroniony (JWT) | –                                                     |
| `GET`   | `/users/{id}` | Pobiera dane użytkownika po UUID. | Chroniony (JWT) | –                                                     |
| `PATCH` | `/users/{id}` | Aktualizuje dane użytkownika (`email`, `name`, `surname`). | Chroniony (JWT) | JSON `{ "email": "", "name": "", "surname": "" }`     |
| `POST`  | `/users/{id}/avatar` | Generuje presigned URL do uploadu avatara. | Chroniony (JWT) | JSON `{ "contentType": "image/png" }`                 |

---

## Przykładowy response dla /auth/me, /users/{id}

```
{
  "id": "b7542757-3b91-4fd0-8e59-473a341f1a3b",
  "email": "user@example.com",
  "name": "Jan",
  "surname": "Kowalski",
  "roles": ["USER"],
  "avatarKey": "avatars/b7542757-3b91-4fd0-8e59-473a341f1a3b/abc123.png"
}
```

## Przykładowy response dla /.well-known/jwks.json

```
{
    "keys": [
        {
            "kty": "RSA",
            "e": "AQAB",
            "kid": "rsa-key-1",
            "n": "hRYrhO4ERL4UUdIEYygve2eamUWI10VNMdCcKvjTJ2-ByLKaoE60EEIcNN1cI79821Qu8gJyiHXcJmYpndHdQXQiXhNHl6HrRDdupEf-tUINkMcFiCAX2tDmMzcxR3D6c8zKf04VbdDtRILoZC51d32vMhKnOT8guahKqBaFIaWcRMukDZfkWjWgWZuy1ITLd4cpCLTfbmZXflOoYNVAZBBlYxbFKcSR4DKxsntJZMS38TDP-tzUKJyy8ksqaEKtaD1-_SswewrWTjd4R5AimPwrtY57ANgtD8vWiTMw_KY66NZoB_Feqtpn2wy-Rhc5KLYRBgi4x4CR-M2g72CchQ"
        }
    ]
}
```

## Przykładowy response dla /oauth2/token

```
{
  "access_token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJyZWNvcmRpbmctYm90LXNlcnZpY2UiLCJzY29wZSI6WyJST0xFX0JPVCIsInVybjpzeXN0ZW06cmVjb3JkaW5nIl0sImlzcyI6Imh0dHA6Ly91c2VyLXNlcnZpY2U6ODA4MCIsImV4cCI6MTcwNTU5NDIwMH0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
  "token_type": "Bearer",
  "expires_in": 300,
  "scope": "ROLE_BOT urn:system:recording"
}
```

## Przykładowy response dla users/{id}/avatar

```
{
  "uploadUrl": "https://bucket.s3.amazonaws.com/avatars/b7542757-3b91-4fd0-8e59-473a341f1a3b/abc123.png?AWSAccessKeyId=..."
}

```
Opis pól:
- kty (Key Type): Typ algorytmu kryptograficznego. Najczęściej będzie to "RSA".
- kid (Key ID): NAJWAŻNIEJSZE POLE. Jest to unikalny identyfikator klucza. Każdy token JWT, który podpisujemy, ma w swoim nagłówku (header) pole kid, które odpowiada jednemu z kluczy na tej liście.
- e (Exponent) i n (Modulus): Te dwa pola razem tworzą klucz publiczny RSA. Twoja biblioteka do obsługi JWT użyje ich do zbudowania klucza publicznego potrzebnego do weryfikacji.



## Przepływ uwierzytelniania

Serwis ten obsługuje dwa główne przepływy uwierzytelniania, oba bazujące na standardzie OAuth 2.0.

1. **Przepływ Użytkownika (Authorization Code Flow)**

    
  Ten przepływ jest używany, gdy człowiek loguje się do systemu za pośrednictwem aplikacji frontendowej i zewnętrznego dostawcy (np. Google). Jest to proces wieloetapowy, oparty na przekierowaniach.

    Kroki:
    
    Inicjalizacja (Frontend ➔ Auth-Service): Użytkownik na frontendzie klika "Zaloguj przez Google". Frontend przekierowuje przeglądarkę użytkownika do endpointu autoryzacji naszego serwera (/oauth2/authorize?client_id=frontend-app...), prosząc o zalogowanie.
    
    Przekierowanie (Auth-Service ➔ Google): Nasz serwer rozpoznaje, że to logowanie dla użytkownika (.oauth2Login()), i przekierowuje przeglądarkę do strony logowania Google.
    
    Logowanie (Użytkownik ➔ Google): Użytkownik loguje się swoimi danymi w domenie Google.
    
    Kod Zwrotny (Google ➔ Auth-Service): Google odsyła przeglądarkę z powrotem do naszego serwisu (na adres .../login/oauth2/code/google) z tymczasowym kodem autoryzacyjnym.
    
    Wymiana i Stworzenie Tokenu (Auth-Service): Nasze handlery (OAuth2AuthenticationSuccessHandler) przechwytują ten kod. Serwis w tle wymienia go z Google na dane użytkownika, a następnie JwtTokenProvider generuje nasz wewnętrzny token JWT.
    
    Finał (Auth-Service ➔ Frontend): Serwis odsyła przeglądarkę użytkownika z powrotem na adres redirectUri frontendu (np. http://localhost:3000/callback), dołączając nowo stworzony token JWT jako parametr URL (?token=...).


2. **Przepływ Serwisu/Bota (Client Credentials Flow)**
   
  Ten przepływ jest używany, gdy zaufana aplikacja backendowa (jak Twój bot nagrywający) potrzebuje uzyskać token do komunikacji M2M (Machine-to-Machine). Ten proces jest bezpośredni i nie wymaga przeglądarki ani interakcji użytkownika.

    Kroki:
    
    - Żądanie Tokenu (Bot ➔ Auth-Service): Bot wysyła bezpośrednie zapytanie POST na endpoint /oauth2/token.
    
    - Uwierzytelnienie Bota: 
      W nagłówku Authorization żądania bot przesyła swoje clientId i clientSecret (np. recording-bot-service i jego hasło) 
      zakodowane w Basic Auth.
      W ciele żądania (jako x-www-form-urlencoded) przesyła grant_type=client_credentials.
    
    - Walidacja i Wydanie Tokenu (Auth-Service): Serwer sprawdza clientId i clientSecret bota w RegisteredClientRepository. 
      Jeśli są poprawne, natychmiast generuje i podpisuje nowy token JWT.
    
    - Odpowiedź (Auth-Service ➔ Bot): Serwer zwraca odpowiedź HTTP 200 OK z ciałem JSON, zawierającym access_token dla bota.