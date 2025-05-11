import { Inject, Logger } from "@nestjs/common";
import { IQueryHandler, QueryHandler } from "@nestjs/cqrs";
import { ChatDTO } from "src/domain/dto/chat.dto";
import { GetChatQuery } from "src/domain/queries/get-chat.query";
import { CHAT_REPOSITORY, IChatRepository } from "src/domain/repository/chat.repository";

@QueryHandler(GetChatQuery)
export class GetMessages implements IQueryHandler<GetChatQuery> {
    private readonly logger = new Logger(GetMessages.name)

    constructor(
        @Inject(CHAT_REPOSITORY)
        private readonly repo: IChatRepository,

    ) { }

    async execute(query: GetChatQuery): Promise<ChatDTO | string> {
        this.logger.log("Getting chat history...");

        const res = await this.repo.getChat(query.id);
        if (!res) {
            this.logger.log("Getting chat history...");
            return "Failed to get chat history"
        }

        this.logger.log("Successfully got history");
        return res;
    }
}