import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { SendMessageHandler } from 'src/app/commands/send-message.command.handler';
import { MESSAGE_REPOSITORY } from 'src/domain/repository/message.repository';
import { DatabaseModule } from 'src/infrastructure/database/database.module';
import { MessageRepositoryImpl } from 'src/infrastructure/database/repositories/message.repository.impl';
import { ReactOnMessageHandler } from './react-on-message.command.handler';
import { ReplyOnMessageHandler } from './reply-on-message.command.handler';

const handlers = [ReactOnMessageHandler, SendMessageHandler, ReplyOnMessageHandler]

@Module({
    imports: [CqrsModule, DatabaseModule],
    providers: [
        ...handlers,
        {
            provide: MESSAGE_REPOSITORY,
            useClass: MessageRepositoryImpl,
        },
    ],
})
export class CommandHandlerModule { }
