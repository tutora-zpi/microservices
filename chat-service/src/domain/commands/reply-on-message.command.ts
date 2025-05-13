import { IsMongoId, IsNotEmpty, IsString, IsUUID } from "class-validator";
import { SendMessageCommand } from "./send-message.command";

export class ReplyOnMessageCommand extends SendMessageCommand {
    @IsNotEmpty()
    @IsUUID()
    readonly replyToMessageID: string; // parent message id 

    constructor(replyToMessageID: string, content: string, senderID: string, receiverID: string, chatID: string) {
        super(receiverID, senderID, chatID, content);
        this.replyToMessageID = replyToMessageID;
    }
}
