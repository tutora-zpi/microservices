import { ReactMessageOnCommand } from "../commands/react-on-message.command";
import { ReplyOnMessageCommand } from "../commands/reply-on-message.command";
import { SendMessageCommand } from "../commands/send-message.command";
import { MessageDTO } from "../dto/message.dto";

export const MESSAGE_REPOSITORY = 'IMessageRepository';

export interface IMessageRepository {
    saveMessage(message: SendMessageCommand): Promise<MessageDTO | null>;
    reactOnMessage(react: ReactMessageOnCommand): Promise<MessageDTO | null>;
    replyOnMessange(reply: ReplyOnMessageCommand): Promise<MessageDTO | null>;
}
