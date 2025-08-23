import { Inject, Logger } from '@nestjs/common';
import { CommandHandler, ICommandHandler } from '@nestjs/cqrs';
import { DeleteChatCommand } from 'src/domain/commands/delete-chat.command';
import { RecordNotFound } from 'src/domain/exceptions/not-found.exception';
import { CHAT_REPOSITORY, IChatRepository } from 'src/domain/repository/chat.repository';


@CommandHandler(DeleteChatCommand)
export class DeleteChatHandler
    implements ICommandHandler<DeleteChatCommand> {
    private readonly logger = new Logger(DeleteChatCommand.name);

    constructor(
        @Inject(CHAT_REPOSITORY)
        private readonly repo: IChatRepository,
    ) { }

    async execute(command: DeleteChatCommand): Promise<void> {
        this.logger.log(`Deleting on chat with id ${command.chatID}`);

        const deletedChat = await this.repo.delete(command);

        if (!deletedChat) {
            this.logger.error('Failed delete chat');
            throw new RecordNotFound();
        }
    }
}
