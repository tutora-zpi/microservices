import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { DatabaseModule } from './infrastructure/database/database.module';
import { QueryHandlerModule } from './app/queries/query.handler.module';
import { ChatModule } from './infrastructure/ws/chat.module';

@Module({
  imports: [ConfigModule.forRoot({
    envFilePath: ".env.local"
  }),
    DatabaseModule, ChatModule, QueryHandlerModule
  ],
  controllers: [],
  providers: [],
})
export class AppModule { }
