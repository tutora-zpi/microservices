import { DeleteChatCommand } from '../commands/delete-chat.command';
import { ChatDTO } from '../dto/chat.dto';
import { MeetingStartedEvent } from '../events/meeting-started.event';
import { GetChatQuery } from '../queries/get-chat.query';

export const CHAT_REPOSITORY = 'IChatRepository';

export interface IChatRepository {
  get(q: GetChatQuery): Promise<ChatDTO | null>;
  init(event: MeetingStartedEvent): Promise<ChatDTO | null>;
  delete(command: DeleteChatCommand): Promise<ChatDTO | null>;
}
