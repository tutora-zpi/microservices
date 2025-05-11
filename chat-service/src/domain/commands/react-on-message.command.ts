import { ICommand } from '@nestjs/cqrs';
import { IsNotEmpty, IsString, Matches } from 'class-validator';

export class ReactMessageOnCommand implements ICommand {
    @IsNotEmpty()
    @IsString()
    messageID: string;

    @IsNotEmpty()
    @IsString()
    userID: string;

    @IsNotEmpty()
    @IsString()
    @Matches(/^("U+1F44D"|"U+2764"|"U+1F525")$/, { message: 'Unknown symbol' }) // emojis like, heart, fire
    emoji: string;

    @IsNotEmpty()
    @IsString()
    chatID: string;
}
