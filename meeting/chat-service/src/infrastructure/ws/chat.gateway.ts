import { Logger } from '@nestjs/common';
import { CommandBus } from '@nestjs/cqrs';
import { WebSocketGateway, WebSocketServer, OnGatewayConnection, OnGatewayDisconnect, MessageBody, SubscribeMessage } from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { ReactMessageOnCommand } from 'src/domain/commands/react-on-message.command';
import { ReplyOnMessageCommand } from 'src/domain/commands/reply-on-message.command';
import { SendMessageCommand } from 'src/domain/commands/send-message.command';
import { UserTyping } from 'src/domain/ws-event/user-typing';

class ErrorResponse {
    constructor(public error: string, public details?: string) { }
}

@WebSocketGateway({
    cors: {
        origin: "*",
    },
    namespace: "/ws/chat",
    pingInterval: 10000,
    pingTimeout: 5000
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
    handleJoinRoom(@MessageBody() data: { roomId: string }, client: Socket) {
        client.join(data.roomId);
        this.logger.log(`Client ${client.id} joined room ${data.roomId}`);
    }

    @SubscribeMessage('sendMessage')
    async handleSendMessage(@MessageBody() data: SendMessageCommand): Promise<void> {
        this.logger.log('Handling send command:', data);

        const command = new SendMessageCommand(data.receiverID, data.senderID, data.meetingID, data.content);

        try {
            const result = await this.commandBus.execute(command);
            this.server.to(data.meetingID).emit('message', result);
        } catch (error) {
            const err = new ErrorResponse('Failed to send message', error.message);
            this.server.to(data.meetingID).emit('message', err);
        }

    }

    @SubscribeMessage('react')
    async handleReact(@MessageBody() data: ReactMessageOnCommand): Promise<void> {
        this.logger.log('Handling react command:', data);

        const command = new ReactMessageOnCommand(data.messageID, data.userID, data.emoji, data.chatID);

        try {
            const result = await this.commandBus.execute(command);
            this.server.to(data.chatID).emit('message', result);
        } catch (error) {
            const err = new ErrorResponse('Failed to react on message', error.message);
            this.server.to(data.chatID).emit('message', err);
        }

    }

    @SubscribeMessage('reply')
    async handleReply(@MessageBody() data: ReplyOnMessageCommand): Promise<void> {
        this.logger.log('Handling react command:', data);

        const command = new ReplyOnMessageCommand(data.replyToMessageID, data.receiverID, data.senderID, data.meetingID, data.content);

        try {
            const result = await this.commandBus.execute(command);
            this.server.to(data.meetingID).emit('message', result);
        } catch (error) {
            const err = new ErrorResponse('Failed to reply on message', error.message);
            this.server.to(data.meetingID).emit('message', err);
        }
    }

    @SubscribeMessage('userTyping')
    handleTyping(@MessageBody() data: UserTyping, client: Socket): void {
        this.logger.debug(`User ${data.userID} is ${data.isTyping ? 'typing' : 'not typing'} in chat ${data.chatID}`);

        client.to(data.chatID).emit('userTyping', data);
    }
}