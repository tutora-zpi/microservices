![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![RabbitMQ](https://img.shields.io/badge/Rabbitmq-FF6600?style=for-the-badge&logo=rabbitmq&logoColor=white)
![Redis](https://img.shields.io/badge/redis-%23DD0031.svg?style=for-the-badge&logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![Amazon S3](https://img.shields.io/badge/Amazon%20S3-FF9900?style=for-the-badge&logo=amazons3&logoColor=white)
![FFmpeg](https://shields.io/badge/FFmpeg-%23171717.svg?logo=ffmpeg&style=for-the-badge&labelColor=171717&logoColor=5cb85c)

# .tutora - Recording Service

Records meetings, preprocesses recordings and uploads them on AWS.

---

## Table of Contents

- [Getting Started](#getting-started)
- [About Recordings ](#about-recordings)

## Getting Started

Copy and fill `.env.local` or `.env` file (depends from running app, local file used to run locally).

```sh
cp .env.sample .env
```
Ensure that you have AWS credintials in your `~/.aws/credintials` and configured bucket as well.

## About Recordings

When the bot receives a request to start recording a room, it connects to the participants, who then start sending audio packets. The bot receives the audio streams and saves them as **.ogg** files.
When the recording stops, all recorded voice tracks (padded with silence where necessary) are merged into a single file.
All resulting artifacts are uploaded to an S3 bucket, and the name of the merged audio file is stored in the session metadata.
After the upload is complete, a [recordings_uploaded_event](/internal/domain/event/recordings_uploaded_event.go) is generated.
