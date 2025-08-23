import { IsNotEmpty, IsUUID } from 'class-validator';
import { SocketEvent } from './socket.event';

export class JoinToRoomSocketEvent extends SocketEvent {
  @IsUUID()
  @IsNotEmpty()
  readonly roomID: string;

  @IsNotEmpty()
  readonly token: string;
}
