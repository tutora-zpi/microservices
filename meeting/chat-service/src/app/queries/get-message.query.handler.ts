import { Inject, Logger, NotFoundException } from "@nestjs/common";
import { IQueryHandler, QueryHandler } from "@nestjs/cqrs";
import { ChatDTO } from "src/domain/dto/chat.dto";
import { GetChatQuery } from "src/domain/queries/get-chat.query";
import { CHAT_REPOSITORY, IChatRepository } from "src/domain/repository/chat.repository";

@QueryHandler(GetChatQuery)
export class GetChatHandler implements IQueryHandler<GetChatQuery> {
    private readonly logger = new Logger(GetChatHandler.name)

    constructor(
        @Inject(CHAT_REPOSITORY)
        private readonly repo: IChatRepository,

    ) { }

    async execute(query: GetChatQuery): Promise<ChatDTO> {
        try {
            this.logger.log("Getting chat history...");

            const res = await this.repo.getChat(query);

            if (!res) {
                const msg = "Chat not found or failed to retrieve chat history";
                this.logger.warn(msg);
                throw new NotFoundException(msg);
            }

            this.logger.log("Successfully retrieved chat history");
            return res;

        } catch (error) {
            if (error instanceof NotFoundException) {
                this.logger.warn(`Chat not found: ${error.message}`);
            } else {
                this.logger.error(`Error retrieving chat: ${error.message}`, error.stack);
            }

            throw error;
        }
    }
}