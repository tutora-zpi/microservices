import { IEvent } from "./event";

import { Type } from "class-transformer";
import { IsUUID, IsNotEmpty, IsArray, ArrayMinSize, ValidateNested } from "class-validator";
import { UserDTO } from "../dto/user.dto";
import { StartMeetingDTO } from "../dto/start-meeting.dto";

export class MeetingStartedEvent implements IEvent {
    @IsUUID()
    @IsNotEmpty()
    readonly meetingID: string;

    @IsArray()
    @ArrayMinSize(2)
    @ValidateNested({ each: true })
    @Type(() => UserDTO)
    @IsNotEmpty({ each: true })
    readonly members: UserDTO[];

    @IsNotEmpty()
    @Type(() => Date)
    readonly startedTime: Date;


    constructor(dto: StartMeetingDTO) {
        this.meetingID = crypto.randomUUID();
        this.members = dto.members;
        this.startedTime = new Date();
    }
}