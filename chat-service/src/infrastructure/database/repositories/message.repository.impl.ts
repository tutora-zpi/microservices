import { Inject, Injectable } from "@nestjs/common";
import { Model } from "mongoose";
import { ReactMessageOnCommand } from "src/domain/commands/react-on-message.command";
import { ReplyOnMessageCommand } from "src/domain/commands/reply-on-message.command";
import { SendMessageCommand } from "src/domain/commands/send-message.command";
import { MessageDTO } from "src/domain/dto/message.dto";
import { Message, MESSAGE_MODEL } from "src/domain/models/message.model";
import { IMessageRepository } from "src/domain/repository/message.repository";


@Injectable()
export class MessageRepositoryImpl implements IMessageRepository {
    constructor(
        @Inject(MESSAGE_MODEL) private readonly collection: Model<Message>,
    ) {
    }

    async replyOnMessange(reply: ReplyOnMessageCommand): Promise<MessageDTO | null> {
        throw new Error("Method not implemented.");
    }

    async reactOnMessage(react: ReactMessageOnCommand): Promise<MessageDTO | null> {
        throw new Error("Method not implemented.");
    }

    async saveMessage(message: SendMessageCommand): Promise<MessageDTO | null> {
        throw new Error("Method not implemented.");
    }

}