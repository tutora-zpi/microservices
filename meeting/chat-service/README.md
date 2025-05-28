pomsyl na czat 

messange trafia na websocket -> wrapujemy to w send message command -> zapis w bazie w command handlerze -> emtiujemy event ze wiadonosc utworzona -> inne serwisy se to obsluguja


jak wyglada czat

start spotkania:    
- init chatu

nie mozna wyslac mess przed powstaniem czatu. 




```json
{
  "pattern": "MeetingStartedEvent",
  "data": {
    "meetingID": "4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a0f",
    "members": [
      {
        "id": "9b7d4d3e-2c36-419c-b90c-d51a5f038bce",
        "avatarURL": "https://example.com/avatar1.png",
        "firstName": "John",
        "lastName": "Doe"
      },
      {
        "id": "f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1c",
        "avatarURL": "https://example.com/avatar2.png",
        "firstName": "Jane",
        "lastName": "Doe"
      }
    ],
    "startedTime": "2025-05-11T20:26:57.023Z"
  }
}
```

# przykladowa odpowiedz na find chat /chats/4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a0f

```json
{
    "data": {
        "id": "4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a0f",
        "members": [
            {
                "id": "9b7d4d3e-2c36-419c-b90c-d51a5f038bce",
                "avatarURL": "https://example.com/avatar1.png",
                "firstName": "John",
                "lastName": "Doe"
            },
            {
                "id": "f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1c",
                "avatarURL": "https://example.com/avatar2.png",
                "firstName": "Jane",
                "lastName": "Doe"
            }
        ],
        "messages": [
            {
                "id": "68359213b04a03c9b239ba3c",
                "chatID": "4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a0f",
                "content": "haha",
                "sender": "f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1c",
                "receiver": "9b7d4d3e-2c36-419c-b90c-d51a5f038bce",
                "reactions": [],
                "answers": [],
                "isRead": false,
                "sentAt": "2025-05-27T10:21:07.716Z"
            }
        ]
    },
    "success": true
}
```