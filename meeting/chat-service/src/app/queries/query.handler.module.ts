import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { GetChatHandler } from 'src/app/queries/get-message.query.handler';
import { CHAT_REPOSITORY } from 'src/domain/repository/chat.repository';
import { DatabaseModule } from 'src/infrastructure/database/database.module';
import { ChatRepositoryImpl } from 'src/infrastructure/database/repositories/chat.repository.impl';

const handlers = [GetChatHandler];

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
