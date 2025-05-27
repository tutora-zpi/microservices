import { Inject, Logger } from "@nestjs/common";
import { CommandHandler, ICommandHandler } from "@nestjs/cqrs";
import { ReplyOnMessageCommand } from "src/domain/commands/reply-on-message.command";

import { MessageDTO } from "src/domain/dto/message.dto";
import { IMessageRepository, MESSAGE_REPOSITORY } from "src/domain/repository/message.repository";

@CommandHandler(ReplyOnMessageCommand)
export class ReplyOnMessageHandler implements ICommandHandler<ReplyOnMessageCommand> {
    private readonly logger = new Logger(ReplyOnMessageHandler.name);

    constructor(
        @Inject(MESSAGE_REPOSITORY)
        private readonly repo: IMessageRepository,
    ) {
    }

    async execute(command: ReplyOnMessageCommand): Promise<MessageDTO | string> {
        const newMessage = await this.repo.reply(command);

        if (!newMessage) {
            this.logger.log("Failed to reply");
            return "Failed to save message";
        }

        // call other services

        this.logger.log("Command successfully excectued")
        return newMessage;
    }
}