import { Inject, Injectable, Logger } from "@nestjs/common";
import { Error, Model } from "mongoose";
import { DeleteMessageCommand } from "src/domain/commands/delete-message.command";
import { ReactMessageOnCommand } from "src/domain/commands/react-on-message.command";
import { ReplyOnMessageCommand } from "src/domain/commands/reply-on-message.command";
import { SendMessageCommand } from "src/domain/commands/send-message.command";
import { MessageDTO } from "src/domain/dto/message.dto";
import { Message, MESSAGE_MODEL } from "src/domain/models/message.model";
import { IMessageRepository } from "src/domain/repository/message.repository";


@Injectable()
export class MessageRepositoryImpl implements IMessageRepository {
    private readonly logger = new Logger(MessageRepositoryImpl.name);

    constructor(
        @Inject(MESSAGE_MODEL) private readonly collection: Model<Message>,
    ) {
    }

    async delete(body: DeleteMessageCommand): Promise<boolean> {
        const res = await this.collection.deleteOne({ _id: body.messageID }).exec();
        this.logger.log("Rows affected:", res.deletedCount)
        return res.deletedCount === 1;
    }

    async reply(reply: ReplyOnMessageCommand): Promise<MessageDTO | null> {
        //TODO: implement        
        throw new Error("Method not implemented.");
    }

    async react(react: ReactMessageOnCommand): Promise<MessageDTO | null> {
        //TODO: implement
        throw new Error("Method not implemented.");
    }

    async save(message: SendMessageCommand): Promise<MessageDTO | null> {
        //TODO: implement
        throw new Error("Method not implemented.");
    }
}