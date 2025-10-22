![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

# .tutora - WS Gateway

Main websocket service used as a layer above some services which demand real-time communication.

---

## Table of Contents

- [Running the app](#running-the-app)
- [Environment Variables](#environment-variables)
- [About service](#about-service)
- [WebSocket Documentation](#ws-documentation)

## Running the app

Service demands `RabbitMQ` and `Redis` so make sure its running in background.

## Environment Variables

By default service uses `.env.local` when its run locally. Use `.env` to provide variables when service is run in container.

## About service

### Why does it even exists?

The main service's purpose is to:

- reduce the number of WebSocket connections on the client side to a single, stable connection,

- extract and centralize shared real-time logic from other services into a single gateway.

### Workflow

Client connects with server and it adds him to connected clients after successfully validation. 

Use case of [UserJoinedHandler](/internal/app/socket_event_handler/general/user_joined_handler.go):

1. Add client to room and retrieve users id.
2. Emit message to room that new user joined.
3. Deliver to client last events from room and board snapshot.

_2 & 3 works concurrently because making snapshot maybe longer_


### Event buffer

This mechanism stores the most recent events and flushes new ones to the message queue at defined intervals.
The goal is to prevent consumers from being overloaded with a high event throughput.

In this context, each event is treated as a package containing both destination metadata and event payload.
This structure allows the broker to precisely determine where each event should be delivered.

### State

When a user joins a room, they initially have no knowledge of the current state. To address this, a simple mechanism delivers the latest content from the cache.

Some events, such as `BoardUpdate`, generate snapshots at specific intervals. Other event types do not rely on continuous updates, so their data is temporarily stored in the cache as well for quick retrieval.

Each snapshot is compressed before being stored and decompressed when fetched from the cache.

## WebSocket Documentation

All handled websocket events you can find [here](/internal/domain/ws_event/).

### Recommendation

It's highly recommended to create some kind of dispatcher on client to easily handle (out/in)coming events.

**Note**: use the native [WebSocket](https://developer.mozilla.org/en-US/docs/Web/API/WebSocket) instead of [Socket.io](https://socket.io/) because service has own [wrapper](/internal/domain/ws_event/ws_event.go) to ws-event.