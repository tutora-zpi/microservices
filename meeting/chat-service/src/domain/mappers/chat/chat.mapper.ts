import { Document } from "mongoose";
import { ChatDTO } from "../../dto/chat.dto";
import { Chat } from "../../models/chat.model";
import { IMapper } from "../mapper";
import { MessageMapper } from "../message/message.mapper";
import { MeetingStartedEvent } from "../../events/meeting-started.event";
import { User } from "../../models/user.model";
import { UserMapper } from "../user/user.mapper";

export class ChatMapper implements IMapper<ChatDTO, Chat> {
    private readonly userMapper = new UserMapper();
    private readonly messageMapper = new MessageMapper();

    toDoc(dto: ChatDTO): Partial<Chat> {
        throw new Error("Method not implemented.");
    }

    toDto(doc: Chat & Document): ChatDTO {
        return {
            id: doc.id,
            members: doc.members.map(user => this.userMapper.toDto(user)),
            messages: doc.messages.map(message => this.messageMapper.toDto(message)),
        };
    }

    fromEvent(event: MeetingStartedEvent): Partial<Chat> {
        return {
            id: event.meetingID,
            members: event.members.map(user => this.userMapper.toDoc(user) as User)
        }
    }
}