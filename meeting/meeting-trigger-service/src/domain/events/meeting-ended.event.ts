import { Type } from "class-transformer";
import { IsUUID, IsNotEmpty, IsArray, ArrayMinSize, ValidateNested } from "class-validator";
import { UserDTO } from "../dto/user.dto";
import { IEvent } from "./event";
import { EndMeetingDTO } from "../dto/end-meeting.dto";

export class MeetingEndedEvent implements IEvent {
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
    readonly endedTime: Date;

    constructor(dto: EndMeetingDTO) {
        this.meetingID = dto.meetingID;
        this.members = dto.members;
        this.endedTime = new Date();
    }
}

