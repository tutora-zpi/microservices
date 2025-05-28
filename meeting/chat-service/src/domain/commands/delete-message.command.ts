import { ICommand } from "@nestjs/cqrs";
import { IsNotEmpty, IsUUID } from "class-validator";

export class DeleteMessageCommand implements ICommand {
    @IsUUID()
    @IsNotEmpty()
    readonly messageID: string;

    constructor(messageID: string) {
        this.messageID = messageID;
    }
}