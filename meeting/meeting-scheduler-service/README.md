![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

# .tutora - MeetingSchedulerService

Responsible for starting/ending and scheduling meetings.  

---


## Table of Contents

- [Running the app](#running-the-app)
- [Environment Variables](#environment-variables)
- [Scheduling](#scheduling)
- [API Documentation](#api-documentation)

## Running the app

Service demands `RabbitMQ` and `Redis` so make sure its running in background.

```bash
docker build -t meeting-scheduler .
docker run -p 8006:8006 --env-file .env.local meeting-scheduler
```

## Environment Variables

Copy **environment variables** exmaple file and fill it with valid data.

```bash
cp .env.sample .env
```
## API Documentation

All documentation you can find on swagger **/api/v1/docs**.

### MeetingDTO:

```ts
interface MeetingDTO {
    meetingId: string
    title: string
    timestamp?: number | null
    members?: UserDTO[] | null
}
```

## Scheduling

Will be provided later...
