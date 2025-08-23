# Directory breakdown

- `account`related services
    - **user-service** responsible for auth user in system.
    - **class-service** class managment service.

- `meeting` related services (mainly real-time)
    - **board-service** handles canva.
    - **voice-service** uses webrtc to enable voice communication.
    - **chat-service** provides functionality to write messages in chat during meeting.

- `ml` related to ai functionalities 
    - **note-service** generating notes from saved .mp3 file.

- `alerts`
    - **notification-service** sending notifications to users.
    - **status-service** tracking users activity in app. 


Nice to have: every service will have `.env.local` file to easily run it locally.
