import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { ChatGateway } from './chat.gateway';
import { MESSAGE_REPOSITORY } from 'src/domain/repository/message.repository';
import { MessageRepositoryImpl } from '../database/repositories/message.repository.impl';
import { DatabaseModule } from '../database/database.module';
import { SecurityModule } from '../security/security.module';
import { CommandHandlerModule } from 'src/app/commands/command.handler.module';

@Module({
  imports: [CqrsModule, DatabaseModule, SecurityModule, CommandHandlerModule],
  providers: [
    ChatGateway,
    { provide: MESSAGE_REPOSITORY, useClass: MessageRepositoryImpl },
  ],
})
export class ChatModule { }
