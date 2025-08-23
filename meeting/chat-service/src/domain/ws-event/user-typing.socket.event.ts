import { IsUUID, IsNotEmpty, IsBoolean } from 'class-validator';
import { SocketEvent } from './socket.event';

export class UserTyping implements SocketEvent {
  @IsUUID()
  @IsNotEmpty()
  readonly chatID: string;

  @IsUUID()
  @IsNotEmpty()
  readonly userID: string;

  @IsBoolean()
  readonly isTyping: boolean;

  constructor(chatID: string, userID: string, isTyping: boolean) {
    this.chatID = chatID;
    this.isTyping = isTyping;
    this.userID = userID;
  }
}
