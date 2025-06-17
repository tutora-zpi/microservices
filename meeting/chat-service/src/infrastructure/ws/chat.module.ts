import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { ChatGateway } from './chat.gateway';
import { MESSAGE_REPOSITORY } from 'src/domain/repository/message.repository';
import { MessageRepositoryImpl } from '../database/repositories/message.repository.impl';
import { SendMessageHandler } from 'src/app/commands/send-message.command.handler';
import { DatabaseModule } from '../database/database.module';
import { ReactMessageOnCommand } from 'src/domain/commands/react-on-message.command';
import { ReplyOnMessageCommand } from 'src/domain/commands/reply-on-message.command';
import { SecurityModule } from '../security/security.module';

const handlers = [SendMessageHandler, ReplyOnMessageCommand, ReactMessageOnCommand,];

@Module({
    imports: [CqrsModule, DatabaseModule, SecurityModule],
    providers: [
        ChatGateway,
        ...handlers,
        {
            provide: MESSAGE_REPOSITORY,
            useClass: MessageRepositoryImpl
        }
    ],
})
export class ChatModule { }