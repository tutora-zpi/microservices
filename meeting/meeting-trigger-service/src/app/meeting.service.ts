import { Inject, Injectable } from '@nestjs/common';
import { ClientProxy } from '@nestjs/microservices';
import { IMeetingService } from './meeting.service.interface';
import { DTO } from 'src/domain/dto/dto';

@Injectable()
export class MeetingService implements IMeetingService {

    constructor(
        @Inject('RABBITMQ_SERVICE') private readonly client: ClientProxy,
    ) { }

    start<T extends DTO>(dto: T): Promise<void> {
        throw new Error('Method not implemented.');
    }
    end<T extends DTO>(dto: T): Promise<void> {
        throw new Error('Method not implemented.');
    }
}
