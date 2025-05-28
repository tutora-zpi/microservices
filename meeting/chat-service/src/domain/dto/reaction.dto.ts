import { IsString, IsNotEmpty, IsUUID } from "class-validator";
import { DTO } from "./dto";
import { UserDTO } from "./user.dto";

export class ReactionDTO extends DTO {
    @IsUUID()
    readonly id: string;

    @IsNotEmpty()
    readonly user: UserDTO;

    @IsString()
    @IsNotEmpty()
    readonly emoji: string;
}
