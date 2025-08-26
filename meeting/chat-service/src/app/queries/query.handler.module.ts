import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { GetChatHandler } from 'src/app/queries/get-chat.query.handler';
import { CHAT_REPOSITORY } from 'src/domain/repository/chat.repository';
import { DatabaseModule } from 'src/infrastructure/database/database.module';
import { ChatRepositoryImpl } from 'src/infrastructure/database/repositories/chat.repository.impl';
import { GetMoreMessagesHandler } from './get-more-messages.handler';

const handlers = [GetChatHandler, GetMoreMessagesHandler];

@Module({
  imports: [CqrsModule, DatabaseModule],
  providers: [
    ...handlers,
    {
      provide: CHAT_REPOSITORY,
      useClass: ChatRepositoryImpl,
    },
  ],
  exports: [...handlers]
})
export class QueryHandlerModule { }
