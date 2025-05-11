import { UserDTO } from "./user.dto";

export class MessageDTO {
    id: string;
    sentAt: Date;
    content: string;
    sender: UserDTO;
    receiver?: UserDTO;
    chatId?: string;
    isRead: boolean;
    reacts: string[];
    answers: MessageDTO[];
}
