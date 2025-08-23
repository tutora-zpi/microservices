import { Module } from '@nestjs/common';
import { CHAT_REPOSITORY } from 'src/domain/repository/chat.repository';
import { DatabaseModule } from 'src/infrastructure/database/database.module';
import { ChatRepositoryImpl } from 'src/infrastructure/database/repositories/chat.repository.impl';
import { MeetingStartedHandler } from './meeting-started.event.handler';

@Module({
  imports: [DatabaseModule],
  providers: [
    {
      provide: CHAT_REPOSITORY,
      useClass: ChatRepositoryImpl,
    },
  ],
  controllers: [MeetingStartedHandler],
})
export class EventHandlerModule {}
