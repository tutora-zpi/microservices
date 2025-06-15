# Serwis Autoryzacji (Auth-Service)

Centralny mikroserwis odpowiedzialny za uwierzytelnianie użytkowników w ekosystemie. Obsługuje logowanie przez dostawców zewnętrznych (OAuth2) i wydaje tokeny JWT podpisane asymetrycznie (RSA), które mogą być weryfikowane przez inne serwisy bez potrzeby odpytywania serwisu autoryzacji.

---

## Funkcjonalności

* **Logowanie jednokrotne (SSO)** z dostawcami:
    * Google
    * GitHub
* **Generowanie tokenów JWT** podpisanych asymetrycznie algorytmem **RS256**.
* **Trwałość użytkowników** w bazie danych PostgreSQL.
* **Endpoint JWKS** (`/.well-known/jwks.json`) do dystrybucji klucza publicznego.
* **Zdecentralizowana walidacja tokenów** – każdy mikroserwis może weryfikować tokeny samodzielnie.
* **Globalna obsługa wyjątków** i spójne formatowanie odpowiedzi błędów.
* Gotowa konfiguracja do uruchomienia za pomocą **Docker Compose**.

## Architektura

![Diagram przepływu OAuth2 i JWT](https://i.imgur.com/gK2gE2C.png)

Serwis implementuje nowoczesny i bezpieczny wzorzec uwierzytelniania dla systemów rozproszonych. Po pomyślnej autoryzacji u dostawcy OAuth2, serwis tworzy lub aktualizuje użytkownika w swojej bazie danych, a następnie generuje token JWT podpisany kluczem prywatnym. Klucz publiczny jest udostępniany przez standardowy endpoint JWKS, co pozwala innym serwisom na bezstanową i szybką weryfikację tokenów.

---

## Wymagania

Przed uruchomieniem upewnij się, że masz zainstalowane następujące narzędzia:
* Java 17 (lub nowsza)
* Apache Maven
* Docker i Docker Compose
* `openssl` (do generowania kluczy RSA, zazwyczaj dostępne w systemach Linux/macOS lub przez Git Bash w Windows)

---

## Konfiguracja

### 1. Generowanie Kluczy RSA

Każde środowisko (deweloperskie i produkcyjne) powinno używać własnej, unikalnej pary kluczy.

Otwórz terminal w głównym katalogu projektu i wykonaj poniższe komendy. Stworzą one katalog `keys` z kluczami, który jest ignorowany przez Git.

```bash
mkdir -p src/main/resources/keys

cd src/main/resources/keys

# Generowanie klucza prywatnego
openssl genpkey -algorithm RSA -out src/main/resources/keys/private_key.pem -pkeyopt rsa_keygen_bits:2048

# Generowanie klucza publicznego
openssl rsa -pubout -in src/main/resources/keys/private_key.pem -out src/main/resources/keys/public_key.pem
```

### 2. Konfiguracja Dostawców OAuth2

Musisz założyć własne aplikacje OAuth2 w konsolach deweloperskich Google i GitHub, aby uzyskać `Client ID` i `Client Secret`.

* **Google:** [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
* **GitHub:** [Developer Settings](https://github.com/settings/developers)

Podczas konfiguracji aplikacji u dostawców, jako **Authorized redirect URI** podaj:
`http://localhost:8080/login/oauth2/code/google`
`http://localhost:8080/login/oauth2/code/github`

### 4. Utworzenie pliku `.env`


Stwórz plik `.env`, skopiuj i uzupełnij poniższy template:
```
# Porty
USER_SERVICE_PORT=
POSTGRES_PORT=
PGADMIN_PORT=

# Dane dostepu do Postgresa
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=

# Dane do pgAdmin
PGADMIN_DEFAULT_EMAIL=
PGADMIN_DEFAULT_PASSWORD=

# OAuth2 (GitHub)
OAUTH_GITHUB_CLIENT_ID=
OAUTH_GITHUB_CLIENT_SECRET=

# OAuth2 (Google) 
OAUTH_GOOGLE_CLIENT_ID=
OAUTH_GOOGLE_CLIENT_SECRET=
```
Następnie **otwórz plik `.env` i uzupełnij go** swoimi danymi uzyskanymi w poprzednim kroku oraz danymi do bazy danych.

**WAŻNE:** Plik `.env` zawiera poufne dane i **nie powinien** być dodawany do systemu kontroli wersji Git. Jest on już uwzględniony w pliku `.gitignore`.

---

## Uruchomienie

Upewnij się, że jesteś w głównym katalogu projektu i wykonaj komendę:
```bash
docker-compose up --build
```


---

## Endpointy API

| Metoda | Ścieżka | Opis | Dostęp |
| :--- | :--- | :--- |:--- |
| `GET` | `/oauth2/authorization/{provider}` | Inicjuje proces logowania przez danego dostawcę (`google` lub `github`). Należy przekierować tu użytkownika. | Publiczny |
| `GET` | `/.well-known/jwks.json` | Zwraca klucz publiczny RSA w formacie JWKS do weryfikacji tokenów. | Publiczny |
| `GET` | `/auth/me` | Zwraca dane o aktualnie zalogowanym użytkowniku. Wymaga nagłówka `Authorization: Bearer <JWT>`. | Chroniony |

---