import { IQuery } from "@nestjs/cqrs";

export class GetMoreMessagesQuery implements IQuery {
    constructor(
        public readonly id: string,
        public readonly page: number,
        public readonly limit: number
    ) {
    }
}