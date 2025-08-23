import { Controller, Inject, Logger } from '@nestjs/common';
import { EventPattern, Payload } from '@nestjs/microservices';
import { MeetingStartedEvent } from 'src/domain/events/meeting-started.event';
import { IEventHandler } from './event.handler.interface';
import {
  CHAT_REPOSITORY,
  IChatRepository,
} from 'src/domain/repository/chat.repository';
import { ChatDTO } from 'src/domain/dto/chat.dto';
import { plainToInstance } from 'class-transformer';
import { validate } from 'class-validator';

@Controller()
export class MeetingStartedHandler
  implements IEventHandler<MeetingStartedEvent> {
  private readonly MAX_RETIRES = 1;

  private readonly logger = new Logger(MeetingStartedHandler.name);

  constructor(
    @Inject(CHAT_REPOSITORY)
    private readonly repo: IChatRepository,
  ) { }

  @EventPattern(MeetingStartedEvent.name)
  async handle(@Payload() event: MeetingStartedEvent) {
    const errors = await validate(plainToInstance(MeetingStartedEvent, event));
    if (errors.length > 0) {
      this.logger.debug('Invalid event');
      return;
    }

    this.logger.log(
      `Received event: ${MeetingStartedEvent.name}`,
      event.meetingID,
    );

    let retries: number = 0;
    let newChat: ChatDTO | null = null;

    try {
      while (retries < this.MAX_RETIRES && newChat == null) {
        try {
          newChat = await this.repo.init(event);
        } catch (err) {
          const msg = err instanceof Error ? err.message : 'Unknown error';
          this.logger.warn('An error occured:', msg);
        }

        if (!newChat) {
          this.logger.log(`Failed to create chat. Retry ${retries + 1}`);
          retries++;
        }
      }

      if (!newChat) {
        // call the meeting-scheduler service to stop in future
        this.logger.log(
          'Error during creating chat. Emitting event to stop meeting.',
        );
        return;
      }

      this.logger.log('Successfully created chat.');
    } catch (error) {
      const msg = error instanceof Error ? error.message : String(error);
      this.logger.debug('Error:', msg);
    }
  }
}
