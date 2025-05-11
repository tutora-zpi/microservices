import { ChatDTO } from "../dto/chat.dto";

export const CHAT_REPOSITORY = 'IChatRepository';

export interface IChatRepository {
    getChat(id: string): Promise<ChatDTO | null>;
    initChat(): Promise<ChatDTO | null>;
}