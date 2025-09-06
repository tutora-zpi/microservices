![NestJS](https://img.shields.io/badge/nestjs-%23E0234E.svg?style=for-the-badge&logo=nestjs&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

# .tutora - ChatService

Responsible for chats and messages in class and during meetings. Real-time communication _enriched_ with answering and reacting to messages.

---

## Table of Contents

- [Running the app](#running-the-app)
- [Environment Variables](#environment-variables)
- [Websocket Documentation](#ws-documentation)
- [API Documentation](#api-documentation)


## Running the app

Before we start ensure you have ran RabbitMQ and MongoDB services from [compose.yml](../../compose.yml). 

**Remember**: if you run **chat-service** in a container, the app will communicate with other containers using the DNS name defined in the file mentioned before.

To run **chat-service** use docker, build image and run it locally.

```bash
docker build -t chat-service .
docker run -p 8002:8002 --env-file .env.local chat-service
```

## Environment Variables

Down below there is a file with env vars, make sure that you have filled it before running an app.   

You can copy content of `.env.sample` and fill it. 

```bash
cp .env.sample .env.local
```

```bash
# App details
APP_NAME=.tutora
PORT=8002

# RabbitMQ passes and ports

RABBITMQ_DEFAULT_USER=
RABBITMQ_DEFAULT_PASS=

RABBITMQ_PORT=5672
RABBITMQ_UI_PORT=15672

# MongoDB passes and ports

MONGO_INITDB_ROOT_USERNAME=
MONGO_INITDB_ROOT_PASSWORD=

MONGO_HOST=
MONGO_PORT=
MONGO_DB_NAME=chat_db

# URLS

RABBITMQ_URL=amqp://<username>:<pass>@<host>:5672/
MONGO_URI=mongodb://<username>:<pass>@<host>:27017/chat_db?authSource=admin

# URL to source of truth - user-service

JWKS_URL=
```

## Websocket Documentation

All events require JWT token.

### Sending message

**Key**: `sendMessage`

**Takes**: [SendMessageCommand](/src/domain/commands/send-message.command.ts)

**Returns**: MessageDTO

**Example**:

```json
{
  "id": "68a98127e0caf5831d37c895",
  "chatID": "room123",
  "content": "asdasd",
  "sender": "f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1d",
  "reactions": [],
  "answers": [],
  "isRead": false,
  "sentAt": "2025-08-23T08:51:51.882Z"
}
```

### React on message

**Key**: `react`

**Takes**: [ReactOnMessageCommand](/src/domain/commands/react-on-message.command.ts)

**Returns**: MessageDTO

**Example**:

```json
{
  "id": "68a980ebe0caf5831d37c861",
  "chatID": "room123",
  "content": "asdasdsad",
  "sender": "9b7d4d3e-2c36-411c-b90c-d51a5f038bce",
  "reactions": [],
  "answers": [
    {
      "id": "68a980efe0caf5831d37c865",
      "chatID": "room123",
      "content": "asdasdsd",
      "sender": {
          "_id": "9b7d4d3e-2c36-411c-b90c-d51a5f038bce",
          "avatarURL": "https://example.com/avatar1.png",
          "firstName": "Johna",
          "lastName": "Doea"
      },
      "reactions": [],
      "answers": [],
      "isRead": false,
      "sentAt": "2025-08-23T08:50:55.084Z"
    }
  ],
  "isRead": false,
  "sentAt": "2025-08-23T08:50:51.410Z"
}
```


### Reply on message

**Key**: `react`

**Takes**: [ReplyOnMessageCommand](/src/domain/commands/reply-on-message.command.ts)

**Returns**: MessageDTO

**Example**:

```json
{
  "id": "68a98253e0caf5831d37c899",
  "chatID": "room123",
  "content": "asdasd",
  "sender": "9b7d4d3e-2c36-411c-b90c-d51a5f038bce",
  "reactions": [
    {
      "id": "68a98255e0caf5831d37c89e",
      "emoji": "❤️",
      "user": {
        "id": "9b7d4d3e-2c36-411c-b90c-d51a5f038bce",
        "avatarURL": "https://example.com/avatar1.png",
        "firstName": "Johna",
        "lastName": "Doea"
      },
      "messageID": "68a98253e0caf5831d37c899"
    }
  ],
  "answers": [],
  "isRead": false,
  "sentAt": "2025-08-23T08:56:51.885Z"
}
```

### User typing

_Note_: ping-pong behaviour.

**Key**: `userTyping`

**Takes**: [UserTypingSocketEvent](/src/domain/ws-event/user-typing.socket.event.ts)

**Returns**: UserTypingSocketEvent


### Joining to the room

**Key**: `joinRoom`

**Takes**: [JoinToRoomSocketEvent](/src/domain/ws-event/join-room.socket.event.ts)

**Action**: adds client to room from payload (meetingID or classID).


## API Documentation

You can find docs on **/api/v1/docs** but down below is additional example responses.

**Note**: every response is wrapped in special response and looks like it:

```json
{
  "success":true, // or false
  "data": {}, // optional dto type
  "error": "", // optional string
}
```

### Creating chat using HTTP (general chat)

**Path**: /api/v1/chats/general

**Exmaple body**: 
```json
{
  "roomID": "4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a0f",
  "members": [
    {
      "id": "9b7d4d3e-2c36-419c-b90c-d51a5f038bce",
      "firstName": "John",
      "lastName": "Doe",
      "avatarURL": "https://example.com/avatar1.png"
    },
    {
      "id": "f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1c",
      "firstName": "Jane",
      "lastName": "Doe",
      "avatarURL": "https://example.com/avatar2.png"
    }
  ]
}
```

**Returns**:
```json
{
  "success": true,
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
    "messages": []
  }
}
```


### Getting more messages

`last_message_id` is optional and should be used when you have last message id after first fetch.

**Path**: /api/v1/chats/{id}/messages?limit={limit}&last_message_id={last_message_id}

**Returns**:

```json
{
  "success": true,
  "data": [
    {
      "id": "68a980f3e0caf5831d37c86d",
      "content": "asdasdsd",
      "sender": "f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1d",
      "reactions": [],
      "answers": [],
      "sentAt": "2025-08-23T08:50:59.626Z"
    },
    {
      "id": "68a98127e0caf5831d37c895",
      "content": "asdasd",
      "sender": "f3ab2bc9-8e22-4f52-b2b4-1b3d73fd6c1d",
      "reactions": [],
      "answers": [],
      "sentAt": "2025-08-23T08:51:51.882Z"
    },
    {
      "id": "68a98253e0caf5831d37c899",
      "content": "asdasd",
      "sender": "9b7d4d3e-2c36-411c-b90c-d51a5f038bce",
      "reactions": [
        {
          "id": "68a98255e0caf5831d37c89e",
          "emoji": "❤️",
          "user": {
            "id": "9b7d4d3e-2c36-411c-b90c-d51a5f038bce",
            "avatarURL": "https://example.com/avatar1.png",
            "firstName": "Johna",
            "lastName": "Doea"
          },
          "messageID": "68a98253e0caf5831d37c899"
        }
      ],
      "answers": [],
      "sentAt": "2025-08-23T08:56:51.885Z"
    }
  ]
}
```

### Finding chat information

**Path**: /api/v1/chats/{id}

**Returns**:

```json
{
  "success": true,
  "data": {
    "id": "4fa0c4f2-3b52-4e61-91a5-bbbd1b2e0a1a",
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
    "messages": []
  }
}
```

### Deleting chat

**Path**: /api/v1/chats/{id}

**Returns**: nothing