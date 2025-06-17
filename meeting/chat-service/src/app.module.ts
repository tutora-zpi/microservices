import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { ClientsModule } from '@nestjs/microservices';
import { DatabaseModule } from './infrastructure/database/database.module';
import { QueryHandlerModule } from './app/queries/query.handler.module';
import { ChatModule } from './infrastructure/ws/chat.module';
import { EventHandlerModule } from './app/events/event.handler.module';
import { ChatControllerModule } from './infrastructure/rest/chat/chat.controller.module';
import { getRmqOptions } from './infrastructure/config/rabbit.config';
import { RabbitMQModule } from './infrastructure/messaging/rabbitmq.module';

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
        useFactory: (configService: ConfigService) => getRmqOptions(configService),
        inject: [ConfigService],
      },
    ]),


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