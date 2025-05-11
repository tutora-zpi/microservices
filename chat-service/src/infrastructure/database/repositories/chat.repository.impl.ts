import { Inject, Injectable } from "@nestjs/common";
import { Model } from "mongoose";
import { ChatDTO } from "src/domain/dto/chat.dto";
import { Chat, CHAT_MODEL } from "src/domain/models/chat.model";
import { IChatRepository } from "src/domain/repository/chat.repository";

@Injectable()
export class ChatRepositoryImpl implements IChatRepository {
    constructor(
        @Inject(CHAT_MODEL) private readonly collection: Model<Chat>,
    ) {
    }

    async getChat(): Promise<ChatDTO | null> {
        throw new Error("Method not implemented.");
    }

    async initChat(): Promise<ChatDTO | null> {
        throw new Error("Method not implemented.");
    }
}