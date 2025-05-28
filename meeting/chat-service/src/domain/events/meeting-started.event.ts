import { Type } from "class-transformer";
import { IsUUID, IsNotEmpty, IsArray, ArrayMinSize, ValidateNested } from "class-validator";
import { UserDTO } from "../dto/user.dto";
import { IEvent } from "@nestjs/cqrs";

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
}
