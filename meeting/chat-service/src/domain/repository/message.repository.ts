import { DeleteMessageCommand } from "../commands/delete-message.command";
import { ReactMessageOnCommand } from "../commands/react-on-message.command";
import { ReplyOnMessageCommand } from "../commands/reply-on-message.command";
import { SendMessageCommand } from "../commands/send-message.command";
import { MessageDTO } from "../dto/message.dto";

export const MESSAGE_REPOSITORY = 'IMessageRepository';

export interface IMessageRepository {
    save(message: SendMessageCommand): Promise<MessageDTO>;
    react(react: ReactMessageOnCommand): Promise<MessageDTO>;
    reply(reply: ReplyOnMessageCommand): Promise<MessageDTO>;
    delete(body: DeleteMessageCommand): Promise<boolean>;
}
