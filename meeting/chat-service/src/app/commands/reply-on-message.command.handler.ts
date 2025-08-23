import { Inject, Logger } from '@nestjs/common';
import { CommandHandler, ICommandHandler } from '@nestjs/cqrs';
import { ReplyOnMessageCommand } from 'src/domain/commands/reply-on-message.command';

import { MessageDTO } from 'src/domain/dto/message.dto';
import {
  IMessageRepository,
  MESSAGE_REPOSITORY,
} from 'src/domain/repository/message.repository';

@CommandHandler(ReplyOnMessageCommand)
export class ReplyOnMessageHandler
  implements ICommandHandler<ReplyOnMessageCommand> {
  private readonly logger = new Logger(ReplyOnMessageHandler.name);

  constructor(
    @Inject(MESSAGE_REPOSITORY)
    private readonly repo: IMessageRepository,
  ) { }

  async execute(command: ReplyOnMessageCommand): Promise<MessageDTO> {
    this.logger.log('Executing command:', command);
    const newMessage = await this.repo.reply(command);

    return newMessage;
  }
}
