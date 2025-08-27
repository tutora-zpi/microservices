import { IsNotEmpty, IsUUID } from 'class-validator';
import { SendMessageCommand } from './send-message.command';

export class ReplyOnMessageCommand extends SendMessageCommand {
  @IsNotEmpty()
  @IsUUID()
  readonly replyToMessageID: string;

  constructor(
    replyToMessageID: string,
    content: string,
    senderID: string,
    chatID: string,
  ) {
    super(senderID, chatID, content);
    this.replyToMessageID = replyToMessageID;
  }
}
