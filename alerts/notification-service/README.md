![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-%234ea94b.svg?style=for-the-badge&logo=mongodb&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

# .tutora - NotificationService

Responsible for sending realtime notifications to clients. Using Server Sent Events (SSE). Additionally server has an endpoint for fetching notfications to show history. 

---


## Table of Contents

- [Running the app](#running-the-app)
- [Environment Variables](#environment-variables)
- [SSE docs](#sse-docs)
- [Events docs](#events-docs)
- [API Documentation](#api-documentation)

## Running the app

Service demands `RabbitMQ` and `MongoDB` so make sure its running in background.

```bash
docker build -t notification-service .
docker run -p 8003:8003 --env-file .env.local notification-service
```

## Environment Variables

Copy **environment variables** exmaple file and fill it with valid data.

```bash
cp .env.sample .env
```

## SSE docs

Down below you'll find example events which are sent from server.

**Note**: good to know that if you handle event to get data simply use `data` property.

If you don't receive notifications in 5 minutes, connection is closed to save resources and client (frontend) should have some kind of reconnect solution.

#### Buffer - What happens if client is reconnecting and in the same time gots message?

If a client is offline and receives a notification, the message is buffered for 30 minutes.

If the buffer reaches its maximum capacity, the oldest messages are removed to make room for new ones.

Even if a notification is removed from the buffer, the client can still access it later by navigating to a dedicated screen or subpage, because all notifications are stored in the database.


#### Connecting

To connect use **/api/v1/notification/stream?token={JWTtoken}**

### Notification Event

**Key**: `notification`

**Data**: server wraps data in [notification dto](/internal/domain/dto/notification_dto.go).

**How can you handle on frontend?**

```js
eventSource.addEventListener("notification", (event) => {
    log(event.data)
    // logic
});
```

### Simple messages

Contains simple messages like heartbeats and initial message.

**Key**: none

**Data**: `string` under data property

**How can you handle on frontend?**

```js
eventSource.onmessage = (event) => {
    log(event.data)
    // logic
}
```

## Events docs

**Note**: every event must be wrapped in specific object.

```json
{
    "pattern": "PatternNameEvent",
    "data": {...}
}
```

Service listens on given types of events. List is down below.

### ClassInvitation Flow

**Note**: if you use **RabbitMQ dashboard** to testing purposes, make sure you paste raw json without comments.

**Events**: [class invitation events](./internal/domain/event/class_invitation/)

`ClassService` publishes an event [ClassInvitationCreatedEvent](./internal/domain/event/class_invitation/created_event.go) then `NotificationService` creates partial entry in database and requests for more data to `UserService` using [UserDetailsRequestedEvent](./internal/domain/event/class_invitation/user_details_requested_event.go). When `UserService` will return data by publishing event called [UserDetailsRespondedEvent](./internal/domain/event/class_invitation/user_details_responded_event.go), `NotificationService` updates fields and creates domain event [ClassInvitationReadyEvent](./internal/domain/event/class_invitation/ready_event.go) and finally pushes notification to infrastructre layer.


### MeetingInivtationEvent

**Events**: [meeting invitation event](./internal/domain/event/meeting_invitation/meeting_invitation_event.go)

Receives an event, saves in database and then pushes notification to infrastructre layer.

## API Documentation

### Fetching notifcations for user

`ReceiverID` is got from **JWT token**.

**Endpoint**: /api/v1/notification?limit={limit}&last_notification_id={id}

`last_notification_id` is optional for the first fetch but next fetches should contain id to easily implement infinite scoll.

**Example**: /api/v1/notification?limit=3&last_notification_id=68b9406b102850c7f8089447

```json
{
    "success": true,
    "data": [
        {
            "id": "68bc910570b5d96a4c46bd54",
            "receiver": {
                "id": "lukasz",
                "firstName": "123",
                "lastName": "456"
            },
            "sender": {
                "id": "tomasz",
                "firstName": "23232",
                "lastName": "3434"
            },
            "createdAt": "2025-09-06T19:52:37Z",
            "type": "invitation",
            "title": "",
            "body": "",
            "redirectionLink": "",
            "metadata": {
                "className": "matma"
            }
        }
    ]
}
```

### Deleting notifcations for user

Used to delete user, returns nothing (204).

**Endpoint**: /api/v1/notification (DELETE)

**Body**:

```json
{
  "ids": ["68bc910570b5d96a4c46bd54", "68bc910570b5d96a4c46bd59"]
}
```

