import { Body, Controller, Post } from "@nestjs/common";
import { IMeetingService } from "src/app/meeting.service.interface";
import { StartMeetingDTO } from "src/domain/dto/start-meeting.dto";

@Controller('meeting')
export class MeetingController {
    constructor(
        private readonly meetingService: IMeetingService
    ) { }


    @Post('start')
    async startMeeting(@Body() meeting: StartMeetingDTO): Promise<void> {
        return await this.meetingService.start(meeting);
    }

    @Post('end')
    async endMeeting(@Body() meeting: StartMeetingDTO): Promise<void> {
        return await this.meetingService.end(meeting);
    }
}
