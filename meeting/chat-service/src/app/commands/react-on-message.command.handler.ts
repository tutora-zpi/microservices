import { Inject, Logger } from '@nestjs/common';
import { CommandHandler, ICommandHandler } from '@nestjs/cqrs';
import { ReactMessageOnCommand } from 'src/domain/commands/react-on-message.command';
import { MessageDTO } from 'src/domain/dto/message.dto';
import {
  IMessageRepository,
  MESSAGE_REPOSITORY,
} from 'src/domain/repository/message.repository';

@CommandHandler(ReactMessageOnCommand)
export class ReactOnMessageHandler
  implements ICommandHandler<ReactMessageOnCommand> {
  private readonly logger = new Logger(ReactMessageOnCommand.name);

  constructor(
    @Inject(MESSAGE_REPOSITORY)
    private readonly repo: IMessageRepository,
  ) { }

  async execute(command: ReactMessageOnCommand): Promise<MessageDTO> {
    this.logger.log(`Reacting on message with id ${command.messageID}`)
    const newMessage = await this.repo.react(command);

    return newMessage;
  }
}
