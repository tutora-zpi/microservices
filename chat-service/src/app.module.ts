import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { ClientsModule, Transport } from '@nestjs/microservices';
import { DatabaseModule } from './infrastructure/database/database.module';
import { QueryHandlerModule } from './app/queries/query.handler.module';
import { ChatModule } from './infrastructure/ws/chat.module';
import { EventHandlerModule } from './app/events/event.handler.module';
import { MeetingStartedEvent } from './domain/events/meeting-started.event';
import { ChatController } from './infrastructure/rest/chat/chat.controller';
import { ChatControllerModule } from './infrastructure/rest/chat/chat.controller.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: ".env.local",
      isGlobal: true,
    }),

    ClientsModule.registerAsync([
      {
        name: 'RABBITMQ_SERVICE',
        imports: [ConfigModule],
        useFactory: (configService: ConfigService) => ({
          transport: Transport.RMQ,
          options: {
            urls: [configService.get<string>('RABBITMQ_URL') || 'amqp://user:user@localhost:5672'],
            queue: MeetingStartedEvent.name,
            queueOptions: {
              durable: true,
            },
            noAck: false,
          },
        }),
        inject: [ConfigService],
      },
    ]),

    DatabaseModule,
    ChatModule,
    QueryHandlerModule,
    EventHandlerModule,
    ChatControllerModule,
  ],
  controllers: [],
  providers: [],
})
export class AppModule { }