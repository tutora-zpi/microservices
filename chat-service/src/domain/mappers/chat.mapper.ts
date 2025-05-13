import { Document } from "mongoose";
import { ChatDTO } from "../dto/chat.dto";
import { Chat } from "../models/chat.model";
import { IMapper } from "./mapper";
import { UserMapper } from "./user.mapper";
import { MessageMapper } from "./message.mapper";
import { MeetingStartedEvent } from "../events/meeting-started.event";
import { User } from "../models/user.model";

export class ChatMapper implements IMapper<ChatDTO, Chat> {
    private readonly userMapper = new UserMapper();
    private readonly messageMapper = new MessageMapper();

    toDoc(dto: ChatDTO): Partial<Chat> {
        throw new Error("Method not implemented.");
    }

    toDto(doc: Chat & Document<unknown, any, any, Record<string, any>>): ChatDTO {
        const dto: ChatDTO = {
            id: doc.id,
            members: doc.members.map(user => this.userMapper.toDto(user)),
            messages: [],
        }

        return dto;
    }

    fromEvent(event: MeetingStartedEvent): Partial<Chat> {
        return {
            _id: event.meetingID,
            members: event.members.map(user => this.userMapper.toDoc(user) as User)
        }
    }
}