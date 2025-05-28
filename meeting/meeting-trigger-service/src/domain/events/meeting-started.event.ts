import { IEvent } from "./event";

import { Type } from "class-transformer";
import { IsUUID, IsNotEmpty, IsArray, ArrayMinSize, ValidateNested } from "class-validator";
import { UserDTO } from "../dto/user.dto";

export class MeetingStartedEvent {
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
}