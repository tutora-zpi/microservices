import { Controller, Logger, Inject } from "@nestjs/common";
import { EventPattern, Payload } from "@nestjs/microservices";
import { IEventHandler } from "./event.handler.interface";
import { BoardUpdateEvent } from "src/domain/events/board.update.event";
import { Scheduler } from "src/infrastructure/scheduler/scheduler.interface";



@Controller()
export class BoardUpdateHandler
    implements IEventHandler<BoardUpdateEvent> {
    private readonly logger: Logger = new Logger(BoardUpdateEvent.name);
    constructor(
        @Inject('Scheduler')
        private readonly scheduler: Scheduler,
    ) { }

    @EventPattern(BoardUpdateEvent.name)
    async handle(@Payload() event: BoardUpdateEvent) {
        this.logger.log("Got: ", event);
        this.scheduler.bufferBoard(event.meetingId, event.data);
    }
}