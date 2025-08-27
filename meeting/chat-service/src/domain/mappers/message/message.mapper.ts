import { MessageDTO } from "../../dto/message.dto";
import { Message } from "../../models/message.model";
import { IMapper } from "../mapper";
import { ReactionMapper } from "../reaction/reaction.mapper";
import { Document } from "mongoose";
import { SendMessageCommand } from "../../commands/send-message.command";

export class MessageMapper implements IMapper<MessageDTO, Message> {
    private readonly reactionMapper = new ReactionMapper();

    toDoc(dto: MessageDTO): Partial<Message> {
        throw new Error("Method not implemented.");
    }

    toDto(doc: Message & Document): MessageDTO {
        return {
            id: doc.id,
            chatID: doc.chatID,
            content: doc.content,
            sender: doc.sender as string,
            reactions: doc.reactions ? doc.reactions.map(reaction => this.reactionMapper.toDto(reaction)) : [],
            answers: doc.answers ? doc.answers.map(answer => this.toDto(answer)) : [],
            isRead: doc.isRead,
            sentAt: doc.sentAt,
        }
    }

    fromCommand(command: SendMessageCommand): Partial<Message> {
        return {
            content: command.content,
            chatID: command.meetingID,
            sender: command.senderID,
        }
    }
}


