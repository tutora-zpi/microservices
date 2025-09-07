## Class controller

Zestaw endpointów do zarządzania klasami (`Classroom`).

### 1. Stwórz nową klasę

Tworzy nową klasę. Użytkownik wysyłający żądanie automatycznie staje się jej gospodarzem (hostem).

-   **Endpoint:** `POST /classes`
-   **Autoryzacja:** Wymagana
-   **Request Body:**

    ```json
    {
      "name": "string"
    }
    ```

-   **Przykładowy Request Body:**

    ```json
    {
      "name": "Analiza matematyczna II"
    }
    ```

-   **Response:** `201 Created`
-   **Przykładowy Response Body:**

    ```json
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Analiza matematyczna II",
      "users": [
        {
          "userId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
          "role": "HOST"
        }
      ]
    }
    ```

***

### 2. Pobierz klasy użytkownika

Pobiera listę wszystkich klas, do których należy uwierzytelniony użytkownik.

-   **Endpoint:** `GET /classes`
-   **Autoryzacja:** Wymagana
-   **Response:** `200 OK`
-   **Przykładowy Response Body:**

    ```json
    [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "name": "Analiza matematyczna II",
        "users": [
          {
            "userId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
            "role": "HOST"
          }
        ]
      },
      {
        "id": "987e6543-e21b-12d3-a456-426614174001",
        "name": "Wprowadzenie do algorytmów",
        "users": [
          {
            "userId": "f1e2d3c4-b5a6-7890-1234-567890abcdef",
            "role": "HOST"
          },
          {
            "userId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
            "role": "MEMBER"
          }
        ]
      }
    ]
    ```

***

### 3. Pobierz szczegóły klasy

Pobiera szczegółowe informacje o konkretnej klasie na podstawie jej ID.

-   **Endpoint:** `GET /classes/{id}`
-   **Autoryzacja:** Niewymagana
-   **Path Variables:**
    -   `id` (UUID): ID klasy do pobrania.
-   **Response:** `200 OK`
-   **Przykładowy Response Body:**

    ```json
    {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "Analiza matematyczna II",
      "users": [
        {
          "userId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
          "role": "HOST"
        },
        {
          "userId": "b2c3d4e5-f6a7-8901-2345-67890abcdef1",
          "role": "GUEST"
        }
      ]
    }
    ```

---

## Invitation controller

Endpointy do zarządzania zaproszeniami do klas.

### 1. Pobierz moje zaproszenia

Pobiera listę wszystkich oczekujących zaproszeń dla uwierzytelnionego użytkownika.

-   **Endpoint:** `GET /invitations/me`
-   **Autoryzacja:** Wymagana
-   **Response:** `200 OK`
-   **Przykładowy Response Body:**

    ```json
    [
      {
        "classId": "123e4567-e89b-12d3-a456-426614174000",
        "userId": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
        "status": "INVITED",
        "createdAt": "2025-09-06T16:30:00Z"
      }
    ]
    ```

***

### 2. Pobierz zaproszenia dla klasy

Pobiera listę wszystkich wysłanych zaproszeń dla konkretnej klasy.

-   **Endpoint:** `GET /invitations/classes/{classId}`
-   **Autoryzacja:** Niewymagana
-   **Path Variables:**
    -   `classId` (UUID): ID klasy.
-   **Response:** `200 OK`
-   **Przykładowy Response Body:**

    ```json
    [
      {
        "classId": "123e4567-e89b-12d3-a456-426614174000",
        "userId": "b2c3d4e5-f6a7-8901-2345-67890abcdef1",
        "status": "DECLINED",
        "createdAt": "2025-09-06T16:40:00Z"
      },
      {
        "classId": "123e4567-e89b-12d3-a456-426614174000",
        "userId": "c3d4e5f6-a7b8-9012-3456-7890abcdef12",
        "status": "ACCEPTED",
        "createdAt": "2025-09-05T10:00:00Z"
      }
    ]
    ```

***

### 3. Zaproś użytkownika do klasy

Wysyła zaproszenie do użytkownika, aby dołączył do klasy, jeżeli nie zostało już wysłane.

-   **Endpoint:** `POST /invitations/{classId}/users/{userId}`
-   **Autoryzacja:** Wymagana
-   **Path Variables:**
    -   `classId` (UUID): ID klasy, do której wysyłane jest zaproszenie.
    -   `userId` (UUID): ID użytkownika, który jest zapraszany.
-   **Response:** `201 Created`
-   **Przykładowy Response Body:**

    ```json
    {
      "classId": "123e4567-e89b-12d3-a456-426614174000",
      "userId": "b2c3d4e5-f6a7-8901-2345-67890abcdef1",
      "status": "INVITED",
      "createdAt": "2025-09-06T16:45:00Z"
    }
    ```

***

### 4. Anuluj zaproszenie

Anuluje wysłane zaproszenie dla użytkownika do klasy, jeżeli nie zostało już zaakceptowane lub odrzucone.

-   **Endpoint:** `DELETE /invitations/{classId}/users/{userId}`
-   **Autoryzacja:** Niewymagana
-   **Path Variables:**
    -   `classId` (UUID): ID klasy.
    -   `userId` (UUID): ID użytkownika, którego zaproszenie jest anulowane.
-   **Response:** `204 No Content`

***

### 5. Akceptuj zaproszenie

Uwierzytelniony użytkownik akceptuje zaproszenie do dołączenia do klasy, jeżeli nie zostało już zaakceptowane lub odrzucone.

-   **Endpoint:** `POST /invitations/{classId}/accept`
-   **Autoryzacja:** Wymagana
-   **Path Variables:**
    -   `classId` (UUID): ID klasy, której zaproszenie jest akceptowane.
-   **Response:** `200 OK`

***

### 6. Odrzuć zaproszenie

Uwierzytelniony użytkownik odrzuca zaproszenie do dołączenia do klasy, jeżeli nie zostało już zaakceptowane lub odrzucone.

-   **Endpoint:** `POST /invitations/{classId}/decline`
-   **Autoryzacja:** Wymagana
-   **Path Variables:**
    -   `classId` (UUID): ID klasy, której zaproszenie jest odrzucane.
-   **Response:** `200 OK`
