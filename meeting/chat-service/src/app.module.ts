import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { DatabaseModule } from './infrastructure/database/database.module';
import { QueryHandlerModule } from './app/queries/query.handler.module';
import { ChatModule } from './infrastructure/ws/chat.module';
import { EventHandlerModule } from './app/events/event.handler.module';
import { ChatControllerModule } from './infrastructure/rest/chat/v1/chat.controller.module';
import { RabbitMQModule } from './infrastructure/messaging/rabbitmq.module';

@Module({
  imports: [
    ConfigModule.forRoot({ envFilePath: '.env.local', isGlobal: true }),
    DatabaseModule,
    ChatModule,
    QueryHandlerModule,
    EventHandlerModule,
    ChatControllerModule,
    RabbitMQModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule { }
