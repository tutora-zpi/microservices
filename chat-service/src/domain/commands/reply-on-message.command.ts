import { IsNotEmpty, IsString } from "class-validator";
import { SendMessageCommand } from "./send-message.command";

export class ReplyOnMessageCommand extends SendMessageCommand {
    @IsNotEmpty()
    @IsString()
    replyToMessageID: string; // parent message id 
}
