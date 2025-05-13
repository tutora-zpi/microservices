import { IQuery } from "@nestjs/cqrs";

export class GetChatQuery implements IQuery {
    constructor(public id: string) {

    }
}