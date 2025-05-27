import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { ChatController } from './chat.controller';
import { QueryHandlerModule } from 'src/app/queries/query.handler.module';
@Module({
    imports: [CqrsModule, QueryHandlerModule],
    controllers: [ChatController],
})
export class ChatControllerModule { }
