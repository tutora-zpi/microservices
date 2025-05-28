import { Inject, Injectable, Logger } from "@nestjs/common";
import mongoose, { Error, Model } from "mongoose";
import { DeleteMessageCommand } from "src/domain/commands/delete-message.command";
import { ReactMessageOnCommand } from "src/domain/commands/react-on-message.command";
import { ReplyOnMessageCommand } from "src/domain/commands/reply-on-message.command";
import { SendMessageCommand } from "src/domain/commands/send-message.command";
import { MessageDTO } from "src/domain/dto/message.dto";
import { UnknownException } from "src/domain/exceptions/unknown.exception";
import { FailedToValidate } from "src/domain/exceptions/validation.exception";
import { MessageMapper } from "src/domain/mappers/message/message.mapper";
import { ReactionMapper } from "src/domain/mappers/reaction/reaction.mapper";
import { Chat, CHAT_MODEL } from "src/domain/models/chat.model";
import { Message, MESSAGE_MODEL } from "src/domain/models/message.model";
import { Reaction, REACTION_MODEL } from "src/domain/models/reaction.model";
import { IMessageRepository } from "src/domain/repository/message.repository";


@Injectable()
export class MessageRepositoryImpl implements IMessageRepository {
    private readonly logger = new Logger(MessageRepositoryImpl.name);
    private readonly mapper = new MessageMapper();
    private readonly reactionMapper = new ReactionMapper();

    constructor(
        @Inject(MESSAGE_MODEL) private readonly messageModel: Model<Message>,
        @Inject(CHAT_MODEL) private readonly chatModel: Model<Chat>,
        @Inject(REACTION_MODEL) private readonly reactionModel: Model<Reaction>,
    ) {
    }

    async delete(body: DeleteMessageCommand): Promise<boolean> {
        const res = await this.messageModel.deleteOne({ _id: body.messageID }).exec();
        this.logger.log("Rows affected:", res.deletedCount)
        return res.deletedCount === 1;
    }

    async reply(reply: ReplyOnMessageCommand): Promise<MessageDTO> {
        this.logger.log("Replying to message with command:", reply);

        try {
            const repliedMessage = await this.save(reply);

            // push id to answers
            const res = await this.messageModel.findByIdAndUpdate(
                reply.replyToMessageID,
                { $push: { answers: repliedMessage.id } },
                { new: true }
            ).exec();


            if (!res) {
                throw new SaveError("Failed to update message with reply");
            }

            return repliedMessage;
        } catch (error) {
            this.logger.error("Error while replying to message:", error.message);

            if (error instanceof FailedToValidate || error instanceof CastingError) {
                throw error;
            }

            throw new UnknownException(`Failed to reply on message ${reply.replyToMessageID}`);
        }

    }

    async react(react: ReactMessageOnCommand): Promise<MessageDTO> {
        this.logger.log("Reacting to message with command:", react);

        const reaction = this.reactionMapper.fromCommand(react);

        try {
            const newReaction = await this.reactionModel.create(reaction);
            if (!newReaction._id) {
                this.logger.error("Failed to retrieve new reaction ID");
                throw new RecordNotFound(`Failed to get new reaction ID for message ${react.messageID}`);
            }

            const updated = await this.messageModel.findOneAndUpdate(
                { _id: react.messageID, chatID: react.chatID },
                { $push: { reactions: newReaction._id } },
                { new: true }
            ).exec();

            if (!updated) {
                this.logger.error("Failed to update message with new reaction");
                throw new SaveError(`Failed to update message with ID ${react.messageID} in chat ${react.chatID}`);
            }

            return this.mapper.toDto(updated);

        } catch (error) {
            this.logger.error("Error while reacting to message:", error.message);
            throw new UnknownException(`Failed to react on message ${react}`);
        }
    }

    async save(message: SendMessageCommand): Promise<MessageDTO> {
        this.logger.log("Saving message with command:", message);


        const newMessage = this.mapper.fromCommand(message);
        this.logger.log("Saving message:", newMessage);

        try {
            const added = await this.messageModel.create(newMessage);

            // update fk
            await this.chatModel.updateOne(
                { id: newMessage.chatID },
                { $push: { messages: added._id } }
            );

            return this.mapper.toDto(added);
        } catch (error) {
            this.logger.error('Error while saving message:', error);

            if (error instanceof mongoose.Error.ValidationError) {
                this.logger.error('Validation erorr:', error.message);

                throw new FailedToValidate(`Failed to validate message ${message}`);

            } else if (error instanceof mongoose.Error.CastError) {
                this.logger.error('Problem with casting:', error.message);

                throw new CastingError(`Failed to cast message ${message}`);
            }

            throw new UnknownException(`Failed to save message ${message}`);
        }
    }
}