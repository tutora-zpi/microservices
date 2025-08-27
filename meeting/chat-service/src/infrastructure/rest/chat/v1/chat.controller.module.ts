import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { ChatController } from './chat.controller';
import { QueryHandlerModule } from 'src/app/queries/query.handler.module';
import { SecurityModule } from 'src/infrastructure/security/security.module';
import { CommandHandlerModule } from 'src/app/commands/command.handler.module';


@Module({
  imports: [CqrsModule, CommandHandlerModule, QueryHandlerModule, SecurityModule],
  controllers: [ChatController],
})
export class ChatControllerModule { }
