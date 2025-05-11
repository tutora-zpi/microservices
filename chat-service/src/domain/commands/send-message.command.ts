import { ICommand } from "@nestjs/cqrs";
import { IsNotEmpty, IsString, Length } from "class-validator";

export class SendMessageCommand implements ICommand {
    @IsString()
    @IsNotEmpty()
    @Length(1, 400)
    content: string;

    @IsString()
    @IsNotEmpty()
    senderID: string;

    @IsString()
    @IsNotEmpty()
    receiverID: string;

    constructor(
        receiverID: string,
        senderID: string,
        content: string
    ) {
        this.content = content;
        this.senderID = senderID;
        this.receiverID = receiverID;
    }

}