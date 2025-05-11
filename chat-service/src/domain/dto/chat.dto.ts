import { MessageDTO } from "./message.dto";

export class ChatDTO {
    id: string;
    members: string[];
    messages: MessageDTO[];
}