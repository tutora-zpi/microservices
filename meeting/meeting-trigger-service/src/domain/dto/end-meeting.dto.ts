import { ArrayMinSize, IsArray, IsNotEmpty, ValidateNested } from "class-validator";
import { DTO } from "./dto";
import { Type } from "class-transformer";
import { UserDTO } from "./user.dto";

export class EndMeetingDTO extends DTO {

    @IsNotEmpty()
    readonly meetingID: string;

    @IsArray()
    @ArrayMinSize(2)
    @ValidateNested({ each: true })
    @Type(() => UserDTO)
    @IsNotEmpty({ each: true })
    readonly members: UserDTO[];
}