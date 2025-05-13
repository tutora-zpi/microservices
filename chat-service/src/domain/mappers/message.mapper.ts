import { MessageDTO } from "../dto/message.dto";
import { Message } from "../models/message.model";
import { IMapper } from "./mapper";
import { UserMapper } from "./user.mapper";
import { ReactionMapper } from "./reaction.mapper";
import { Document } from "mongoose";

export class MessageMapper implements IMapper<MessageDTO, Message> {
    private readonly userMapper = new UserMapper();
    private readonly reactionMapper = new ReactionMapper();
    toDoc(dto: MessageDTO): Partial<Message> {
        throw new Error("Method not implemented.");
    }
    toDto(doc: Message & Document<unknown, any, any, Record<string, any>>): MessageDTO {
        throw new Error("Method not implemented.");
    }
}
