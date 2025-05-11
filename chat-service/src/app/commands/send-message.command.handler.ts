import { Inject, Injectable, Logger } from "@nestjs/common";
import { CommandHandler, ICommandHandler } from "@nestjs/cqrs";
import { SendMessageCommand } from "src/domain/commands/send-message.command";
import { MessageDTO } from "src/domain/dto/message.dto";
import { IMessageRepository, MESSAGE_REPOSITORY } from "src/domain/repository/message.repository";

@CommandHandler(SendMessageCommand)
export class SendMessageHandler implements ICommandHandler<SendMessageCommand> {
    private readonly logger = new Logger(SendMessageHandler.name);

    constructor(
        @Inject(MESSAGE_REPOSITORY)
        private readonly repo: IMessageRepository,
    ) {
    }

    async execute(command: SendMessageCommand): Promise<MessageDTO | string> {
        const newMessage = await this.repo.saveMessage(command);

        if (!newMessage) {
            this.logger.log("Failed to save message");
            return "Failed to save message";
        }

        // call other services

        this.logger.log("Command successfully excectued")
        return newMessage;
    }
}