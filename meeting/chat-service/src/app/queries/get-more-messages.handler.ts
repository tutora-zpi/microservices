import { Logger, Inject } from "@nestjs/common";
import { QueryHandler, IQueryHandler } from "@nestjs/cqrs";
import { MessageDTO } from "src/domain/dto/message.dto";
import { RecordNotFound } from "src/domain/exceptions/not-found.exception";
import { GetMoreMessagesQuery } from "src/domain/queries/get-more-messages.query";
import { IMessageRepository, MESSAGE_REPOSITORY } from "src/domain/repository/message.repository";

@QueryHandler(GetMoreMessagesQuery)
export class GetMoreMessagesHandler implements IQueryHandler<GetMoreMessagesQuery> {
    private readonly logger = new Logger(GetMoreMessagesHandler.name);

    constructor(
        @Inject(MESSAGE_REPOSITORY)
        private readonly repo: IMessageRepository,
    ) { }

    async execute(query: GetMoreMessagesQuery): Promise<MessageDTO[]> {
        this.logger.log(`Executing GetChatQuery with id: ${query.id}`);

        const res = await this.repo.get(query);
        if (!res) {
            throw new RecordNotFound();
        }

        return res;
    }
}