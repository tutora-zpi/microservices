import { Module } from '@nestjs/common';
import { MeetingController } from 'src/infrastrucutre/rest/meeting.controller';
import { MeetingService } from './meeting.service';
import { RabbitModule } from 'src/infrastrucutre/rabbit/rabbit.module';

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

