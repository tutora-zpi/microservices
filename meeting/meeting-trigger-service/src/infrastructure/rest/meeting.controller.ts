import { BadRequestException, Body, Controller, HttpCode, Inject, Post, Put, UseGuards } from "@nestjs/common";
import { AuthGuard } from "@nestjs/passport";
import { IMeetingService } from "src/app/meeting.service.interface";
import { EndMeetingDTO } from "src/domain/dto/end-meeting.dto";
import { StartMeetingDTO } from "src/domain/dto/start-meeting.dto";

@Controller('meeting')
export class MeetingController {
    constructor(
        @Inject('IMeetingService')
        private readonly meetingService: IMeetingService
    ) { }


    @UseGuards(AuthGuard('jwt'))
    @Post('start')
    @HttpCode(201)
    async startMeeting(@Body() meeting: StartMeetingDTO): Promise<string> {

        try {
            const id = await this.meetingService.start(meeting);

            return id;
        } catch (error) {
            throw new BadRequestException(`Failed to start meeting: ${error.message}`);
        }
    }

    // maybe change it end/meetingID, we might not need enitre dto  
    @UseGuards(AuthGuard('jwt'))
    @Put('end')
    @HttpCode(200)
    async endMeeting(@Body() meeting: EndMeetingDTO): Promise<void> {
        try {
            await this.meetingService.end(meeting);
        } catch (error) {
            throw new BadRequestException(`Failed to end meeting: ${error.message}`);
        }
    }
}
