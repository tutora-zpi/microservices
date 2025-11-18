import { IEvent } from '@nestjs/cqrs';


export interface IEventHandler<T extends IEvent> {
    handle(event: T): Promise<void>;
}