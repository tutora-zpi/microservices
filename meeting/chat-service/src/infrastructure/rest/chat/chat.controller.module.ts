import { Module } from '@nestjs/common';
import { CqrsModule } from '@nestjs/cqrs';
import { ChatController } from './chat.controller';
import { QueryHandlerModule } from 'src/app/queries/query.handler.module';
import { SecurityModule } from 'src/infrastructure/security/security.module';
@Module({
    imports: [CqrsModule, QueryHandlerModule, SecurityModule],
    controllers: [ChatController],
})
export class ChatControllerModule { }
