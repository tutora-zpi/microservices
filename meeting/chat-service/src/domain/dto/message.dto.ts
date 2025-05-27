import { DTO } from "./dto";
import { ReactionDTO } from "./reaction.dto";
import { UserDTO } from "./user.dto";
import { IsArray, IsBoolean, IsDate, IsNotEmpty, IsOptional, IsString, IsMongoId } from 'class-validator';

export class MessageDTO extends DTO {
    @IsMongoId()
    readonly id: string;

    @IsDate()
    readonly sentAt: Date;

    @IsString()
    @IsNotEmpty()
    readonly content: string;

    readonly sender: UserDTO;

    @IsOptional()
    readonly receiver?: UserDTO;

    @IsMongoId()
    readonly chatID: string;

    @IsBoolean()
    readonly isRead: boolean;

    @IsArray()
    readonly reactions: ReactionDTO[];

    @IsArray()
    readonly answers: MessageDTO[];
}