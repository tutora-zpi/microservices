import { Inject, Logger } from '@nestjs/common';
import { CommandHandler, ICommandHandler } from '@nestjs/cqrs';
import { CreateGeneralChatCommand } from 'src/domain/commands/create-general-chat.command';
import { ChatDTO } from 'src/domain/dto/chat.dto';
import { MeetingStartedEvent } from 'src/domain/events/meeting-started.event';
import { CouldNotCreateGeneralChat } from 'src/domain/exceptions/could-not-create-general-chat.exception';
import { CHAT_REPOSITORY, IChatRepository } from 'src/domain/repository/chat.repository';


@CommandHandler(CreateGeneralChatCommand)
export class CreateGeneralChatHandler
    implements ICommandHandler<CreateGeneralChatCommand> {
    private readonly logger = new Logger(CreateGeneralChatCommand.name);

    constructor(
        @Inject(CHAT_REPOSITORY)
        private readonly repo: IChatRepository,
    ) { }

    async execute(command: CreateGeneralChatCommand): Promise<ChatDTO> {
        this.logger.log('Executing command:', command);

        const transformedToEvent: MeetingStartedEvent = {
            meetingID: command.roomID,
            members: command.members,
        }

        const newChat = await this.repo.init(transformedToEvent);

        if (!newChat) {
            this.logger.error("ChatDTO has not been returned");
            throw new CouldNotCreateGeneralChat();
        }

        return newChat;
    }
}
