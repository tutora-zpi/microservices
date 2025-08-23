import { Logger, UseGuards } from '@nestjs/common';
import { CommandBus } from '@nestjs/cqrs';
import {
  WebSocketGateway,
  WebSocketServer,
  OnGatewayConnection,
  OnGatewayDisconnect,
  MessageBody,
  SubscribeMessage,
  ConnectedSocket,
} from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { ReactMessageOnCommand } from 'src/domain/commands/react-on-message.command';
import { ReplyOnMessageCommand } from 'src/domain/commands/reply-on-message.command';
import { SendMessageCommand } from 'src/domain/commands/send-message.command';
import { JoinToRoomSocketEvent } from 'src/domain/ws-event/join-room.socket.event';
import { UserTyping as UserTypingSocketEvent } from 'src/domain/ws-event/user-typing.socket.event';
import { WsAuthGuard } from '../security/guards/ws.auth.guard';
import { MessageDTO } from 'src/domain/dto/message.dto';

class ErrorResponse {
  constructor(
    public error: string,
    public details?: string,
  ) { }
}

@UseGuards(WsAuthGuard)
@WebSocketGateway({
  cors: {
    origin: '*',
  },
  namespace: '/ws/chat',
  pingInterval: 10000,
  pingTimeout: 5000,
})
export class ChatGateway implements OnGatewayConnection, OnGatewayDisconnect {
  private readonly logger = new Logger(ChatGateway.name);

  constructor(private readonly commandBus: CommandBus) { }

  @WebSocketServer() server: Server;

  handleConnection(client: Socket) {
    this.logger.log('Client connected: ' + client.id);
  }

  handleDisconnect(client: Socket) {
    this.logger.log('Client disconnected: ' + client.id);
  }

  @SubscribeMessage('joinRoom')
  async handleJoinRoom(
    @MessageBody() data: JoinToRoomSocketEvent,
    @ConnectedSocket() client: Socket,
  ) {
    await client.join(data.roomID);
    this.logger.log(`Client ${client.id} joined room ${data.roomID}`);
  }

  @SubscribeMessage('sendMessage')
  async handleSendMessage(
    @MessageBody() data: SendMessageCommand,
  ): Promise<void> {
    this.logger.log('Handling send command:', data);

    const command = new SendMessageCommand(
      data.senderID,
      data.meetingID,
      data.content,
    );

    try {
      const result =
        await this.commandBus.execute<SendMessageCommand, MessageDTO>(command);

      this.logger.debug('Result [handleSendMessage]:', result);

      this.server.to(data.meetingID).emit('message', result);
    } catch (error) {
      const msg = error instanceof Error ? error.message : String(error);

      const err = new ErrorResponse('Failed to send message', msg);

      this.server.to(data.meetingID).emit('message', err);
    }
  }

  @SubscribeMessage('react')
  async handleReact(@MessageBody() data: ReactMessageOnCommand): Promise<void> {
    const command = new ReactMessageOnCommand(
      data.messageID,
      data.userID,
      data.emoji,
      data.chatID,
    );

    try {
      const result =
        await this.commandBus.execute<ReactMessageOnCommand, MessageDTO>(command);

      this.logger.debug('Result [handleReact]: ', result);

      this.server.to(command.chatID).emit('message', result);
    } catch (error) {
      const msg = error instanceof Error ? error.message : String(error);

      const err = new ErrorResponse('Failed to react on message', msg);

      this.server.to(data.chatID).emit('message', err);
    }
  }

  @SubscribeMessage('reply')
  async handleReply(@MessageBody() data: ReplyOnMessageCommand): Promise<void> {
    const command = new ReplyOnMessageCommand(
      data.replyToMessageID,
      data.content,
      data.senderID,
      data.meetingID,
    );

    try {
      const result =
        await this.commandBus.execute<ReplyOnMessageCommand, MessageDTO>(command);

      this.logger.debug('Result [handleReply]: ', result);

      this.server.to(command.meetingID).emit('message', result);
    } catch (error) {
      const msg = error instanceof Error ? error.message : String(error);

      const err = new ErrorResponse('Failed to reply on message', msg);

      this.server.to(command.meetingID).emit('message', err);
    }
  }

  @SubscribeMessage('userTyping')
  handleTyping(
    @MessageBody() data: UserTypingSocketEvent,
    @ConnectedSocket() client: Socket,
  ): void {
    const typingEvent = new UserTypingSocketEvent(
      data.chatID,
      data.userID,
      data.isTyping,
    );

    this.logger.debug(
      `User ${typingEvent.userID} is ${typingEvent.isTyping ? 'typing' : 'stopped typing'} in chat ${data.chatID}`,
    );

    client.to(typingEvent.chatID).emit('userTyping', typingEvent);
  }

}
