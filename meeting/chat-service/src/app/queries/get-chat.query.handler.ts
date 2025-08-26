import { Inject, Logger } from '@nestjs/common';
import { IQueryHandler, QueryHandler } from '@nestjs/cqrs';
import { ChatDTO } from 'src/domain/dto/chat.dto';
import { RecordNotFound } from 'src/domain/exceptions/not-found.exception';
import { GetChatQuery } from 'src/domain/queries/get-chat.query';
import {
  CHAT_REPOSITORY,
  IChatRepository,
} from 'src/domain/repository/chat.repository';

@QueryHandler(GetChatQuery)
export class GetChatHandler implements IQueryHandler<GetChatQuery> {
  private readonly logger = new Logger(GetChatHandler.name);

  constructor(
    @Inject(CHAT_REPOSITORY)
    private readonly repo: IChatRepository,
  ) { }

  async execute(query: GetChatQuery): Promise<ChatDTO> {
    this.logger.log(`Executing GetChatQuery with id: ${query.id}`);

    const res = await this.repo.get(query);
    if (!res) {
      throw new RecordNotFound();
    }

    return res;
  }
}
