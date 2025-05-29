import { Inject, Injectable, Logger } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { IMeetingService } from './meeting.service.interface';
import { MeetingStartedEvent } from 'src/domain/events/meeting-started.event';
import { MeetingEndedEvent } from 'src/domain/events/meeting-ended.event';
import { firstValueFrom } from 'rxjs';
import { StartMeetingDTO } from 'src/domain/dto/start-meeting.dto';
import { EndMeetingDTO } from 'src/domain/dto/end-meeting.dto';
import { FailedToPublishEvent } from 'src/domain/exceptions/publish-event-fail.exception';
import { RABBITMQ_SERVICE } from 'src/infrastructure/rabbit/rabbit.module';

@Injectable()
export class MeetingService implements IMeetingService {
    private readonly logger = new Logger(MeetingService.name);

    constructor(
        @Inject(RABBITMQ_SERVICE) private readonly client: ClientProxy,
    ) { }

    async start(dto: StartMeetingDTO): Promise<void> {
        try {
            await firstValueFrom(this.client.emit<MeetingStartedEvent>(
                MeetingStartedEvent.name,
                new MeetingStartedEvent(dto),
            ));
            this.logger.log('MeetingStartedEvent sent');
        } catch (err) {
            this.logger.error('Failed to emit MeetingStartedEvent', err);
            throw new FailedToPublishEvent("Failed to publish MeetingStartedEvent");
        }
    }

    async end(dto: EndMeetingDTO): Promise<void> {
        try {
            await firstValueFrom(this.client.emit<MeetingEndedEvent>(
                MeetingEndedEvent.name,
                new MeetingEndedEvent(dto)
            ));

            this.logger.log('MeetingEndedEvent sent');
        } catch (err) {
            this.logger.error('Failed to emit MeetingEndedEvent', err);
            throw new FailedToPublishEvent("Failed to publish MeetingEndedEvent");
        }
    }
}
