import { IQuery } from "@nestjs/cqrs";

export class GetMoreMessagesQuery implements IQuery {
    constructor(
        public readonly id: string,
        public readonly limit: number,
        public readonly lastMessageId?: string | null,
    ) {
    }
}