import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { ChatGateway } from './chat.gateway';
import { MESSAGE_REPOSITORY } from 'src/domain/repository/message.repository';
import { MessageRepositoryImpl } from '../database/repositories/message.repository.impl';
import { SendMessageHandler } from 'src/app/commands/send-message.command.handler';
import { DatabaseModule } from '../database/database.module';

const CommandHandlers = [SendMessageHandler];

@Module({
    imports: [CqrsModule, DatabaseModule],
    providers: [
        ChatGateway,
        ...CommandHandlers,
        {
            provide: MESSAGE_REPOSITORY,
            useClass: MessageRepositoryImpl
        }
    ],
})
export class ChatModule { }