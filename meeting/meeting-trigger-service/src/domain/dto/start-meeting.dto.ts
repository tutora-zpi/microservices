import { Type } from "class-transformer";
import { IsArray, ArrayMinSize, ValidateNested, IsNotEmpty } from "class-validator";
import { DTO } from "./dto";
import { UserDTO } from "./user.dto";

export class StartMeetingDTO extends DTO {
    @IsArray()
    @ArrayMinSize(2)
    @ValidateNested({ each: true })
    @Type(() => UserDTO)
    @IsNotEmpty({ each: true })
    readonly members: UserDTO[];
}

