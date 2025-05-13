import { DeleteMessageCommand } from "../commands/delete-message.command";
import { ReactMessageOnCommand } from "../commands/react-on-message.command";
import { ReplyOnMessageCommand } from "../commands/reply-on-message.command";
import { SendMessageCommand } from "../commands/send-message.command";
import { MessageDTO } from "../dto/message.dto";

export const MESSAGE_REPOSITORY = 'IMessageRepository';

export interface IMessageRepository {
    save(message: SendMessageCommand): Promise<MessageDTO | null>;
    react(react: ReactMessageOnCommand): Promise<MessageDTO | null>;
    reply(reply: ReplyOnMessageCommand): Promise<MessageDTO | null>;
    delete(body: DeleteMessageCommand): Promise<boolean>;
}
