import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { SendMessageHandler } from 'src/app/commands/send-message.command.handler';
import { MESSAGE_REPOSITORY } from 'src/domain/repository/message.repository';
import { DatabaseModule } from 'src/infrastructure/database/database.module';
import { MessageRepositoryImpl } from 'src/infrastructure/database/repositories/message.repository.impl';

@Module({
    imports: [CqrsModule, DatabaseModule],
    providers: [
        SendMessageHandler,
        {
            provide: MESSAGE_REPOSITORY,
            useClass: MessageRepositoryImpl,
        },
    ],
})
export class CommandHandlerModule { }
