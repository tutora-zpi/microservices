import { IEvent } from "@nestjs/cqrs";

export class BoardUpdateEvent implements IEvent {
    constructor(
        public readonly meetingId: string,
        public readonly data: any,
    ) { }
}