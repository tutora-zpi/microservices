import { Module } from '@nestjs/common';
import { ConfigModule, } from '@nestjs/config';
import { MeetingModule } from './app/meeting.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      envFilePath: ".env.local",
      isGlobal: true,
    }),
    MeetingModule,
  ],
})

export class AppModule { }
