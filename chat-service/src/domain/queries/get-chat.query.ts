import { IQuery } from "@nestjs/cqrs";

export class GetChatQuery implements IQuery {
    id: string;
}