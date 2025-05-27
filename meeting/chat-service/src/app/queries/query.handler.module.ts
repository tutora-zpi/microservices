import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { GetChatHandler } from 'src/app/queries/get-message.query.handler';
import { CHAT_REPOSITORY } from 'src/domain/repository/chat.repository';
import { DatabaseModule } from 'src/infrastructure/database/database.module';
import { ChatRepositoryImpl } from 'src/infrastructure/database/repositories/chat.repository.impl';

@Module({
    imports: [CqrsModule, DatabaseModule],
    providers: [
        GetChatHandler,
        {
            provide: CHAT_REPOSITORY,
            useClass: ChatRepositoryImpl,
        },
    ],
})
export class QueryHandlerModule { }
