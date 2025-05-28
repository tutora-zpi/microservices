import { ChatDTO } from "../dto/chat.dto";
import { MeetingStartedEvent } from "../events/meeting-started.event";
import { GetChatQuery } from "../queries/get-chat.query";

export const CHAT_REPOSITORY = 'IChatRepository';

export interface IChatRepository {
    getChat(q: GetChatQuery): Promise<ChatDTO>;
    initChat(event: MeetingStartedEvent): Promise<ChatDTO>;
}