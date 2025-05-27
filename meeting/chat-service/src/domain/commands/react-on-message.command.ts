import { ICommand } from '@nestjs/cqrs';
import { IsNotEmpty, IsString, IsUUID, Matches, } from 'class-validator';

export class ReactMessageOnCommand implements ICommand {
    @IsNotEmpty()
    @IsUUID()
    readonly messageID: string;

    @IsNotEmpty()
    @IsUUID()
    readonly userID: string;

    @IsNotEmpty()
    @IsString()
    @Matches(/^("U+1F44D"|"U+2764"|"U+1F525")$/, { message: 'Unknown emoji' }) // emojis like, heart, fire
    readonly emoji: string;

    // probably unused
    @IsNotEmpty()
    @IsUUID()
    readonly chatID: string;

    constructor(messageID: string, userID: string, emoji: string, chatID: string) {
        this.messageID = messageID;
        this.userID = userID;
        this.emoji = emoji;
        this.chatID = chatID;
    }
}
