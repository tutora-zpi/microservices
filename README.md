# Directory breakdown

- `account`related services
    - **user-service** responsible for auth user in system.
    - **class-service** class managment service.

- `meeting` related services (mainly real-time)
    - **meeting-scheduler-service** main service to create and inviting people meeting.
    - **board-service** handles canva.
    - **recorder-service** recording/saving voice during meeting.
    - **chat-service** provides functionality to write messages in chat during meeting.
    - **signaling-service** 

- `ml` related to ai functionalities 
    - **note-service** generating notes from saved .mp3 file.

- `alerts`
    - **notification-service** sending notifications to users.
    - **status-service** tracking users activity in app. 
