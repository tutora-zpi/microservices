import { Module } from '@nestjs/common';
import { MeetingController } from 'src/infrastructure/rest/meeting.controller';
import { MeetingService } from './meeting.service';
import { RabbitModule } from 'src/infrastructure/rabbit/rabbit.module';

@Module({
    imports: [
        RabbitModule,
    ],
    controllers: [MeetingController],
    providers: [
        { provide: 'IMeetingService', useClass: MeetingService },
    ],
})
export class MeetingModule { }

