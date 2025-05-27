import { DTO } from "./dto";
import { MessageDTO } from "./message.dto";
import { UserDTO } from "./user.dto";
import { IsArray, IsDate, IsNotEmpty, IsOptional, IsUUID } from 'class-validator';

export class ChatDTO extends DTO {
    @IsUUID()
    @IsNotEmpty()
    readonly id: string;

    @IsArray()
    readonly members: UserDTO[];

    @IsArray()
    readonly messages: MessageDTO[];

    @IsOptional()
    @IsDate()
    readonly createdAt?: Date;

    @IsOptional()
    @IsDate()
    readonly updatedAt?: Date;


}