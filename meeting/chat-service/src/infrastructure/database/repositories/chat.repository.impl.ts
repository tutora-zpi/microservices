import { Inject, Injectable, Logger } from "@nestjs/common";
import { Model } from "mongoose";
import { ChatDTO } from "src/domain/dto/chat.dto";
import { MeetingStartedEvent } from "src/domain/events/meeting-started.event";
import { UnknownException } from "src/domain/exceptions/unknown.exception";
import { ChatMapper } from "src/domain/mappers/chat/chat.mapper";
import { Chat, CHAT_MODEL } from "src/domain/models/chat.model";
import { Message } from "src/domain/models/message.model";
import { User } from "src/domain/models/user.model";
import { GetChatQuery } from "src/domain/queries/get-chat.query";
import { IChatRepository } from "src/domain/repository/chat.repository";
import { IUserRepository, USER_REPOSITORY } from "src/domain/repository/user.repository";

@Injectable()
export class ChatRepositoryImpl implements IChatRepository {
    private readonly logger: Logger = new Logger(ChatRepositoryImpl.name);
    private readonly mapper: ChatMapper = new ChatMapper(); // in future we can change it to injectable mapper

    constructor(
        @Inject(CHAT_MODEL) private readonly chatModel: Model<Chat>,

        @Inject(USER_REPOSITORY) private readonly userRepo: IUserRepository,
    ) {
    }

    async getChat(q: GetChatQuery): Promise<ChatDTO> {
        try {
            const chat = await this.chatModel.findOne({ id: q.id })
                .populate<User>({
                    path: 'members',
                })
                .populate<Message>({
                    path: 'messages',
                    options: { sort: { createdAt: -1 } },
                })
                .exec();

            if (!chat) {
                this.logger.log("Could not find chat with", q.id);
                throw new RecordNotFound(`Chat with id ${q.id} not found`);
            }

            this.logger.debug(chat);

            const dto = this.mapper.toDto(chat);

            return dto;
        } catch (error) {
            this.logger.error("Failed to retrieve chat", error);
            throw new UnknownException(`Failed to get chat with id ${q.id}`);
        }
    }

    async initChat(event: MeetingStartedEvent): Promise<ChatDTO> {
        try {
            const newChat = this.mapper.fromEvent(event);

            // save users 
            const users = await this.userRepo.saveUsers(newChat.members as User[]);

            if (!users) {
                throw new SaveError("Failed to save users");
            }

            const fks = users.map(u => u.id);

            const chat = await this.chatModel.findOneAndUpdate(
                { id: newChat.id },
                {
                    $set: {
                        id: newChat.id,
                        ...(newChat.messages && { messages: newChat.messages })
                    },
                    $addToSet: {
                        members: { $each: fks }
                    }
                },
                {
                    upsert: true,
                    new: true,
                    setDefaultsOnInsert: true
                }
            ).populate('members').exec();

            return this.mapper.toDto(chat);

        } catch (error) {
            this.logger.error("Something went wrong", error.message);
            throw new UnknownException("Failed to init chat");
        }
    }
}