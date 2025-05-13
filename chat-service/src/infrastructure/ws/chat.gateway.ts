import { Logger } from '@nestjs/common';
import { CommandBus } from '@nestjs/cqrs';
import { WebSocketGateway, WebSocketServer, OnGatewayConnection, OnGatewayDisconnect, MessageBody, SubscribeMessage } from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
import { ReactMessageOnCommand } from 'src/domain/commands/react-on-message.command';
import { ReplyOnMessageCommand } from 'src/domain/commands/reply-on-message.command';
import { SendMessageCommand } from 'src/domain/commands/send-message.command';

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

        const command = new SendMessageCommand(data.receiverID, data.senderID, data.meetingID, data.content);

        const result = await this.commandBus.execute(command);

        this.server.to(data.meetingID).emit('message', result);
    }

    @SubscribeMessage('react')
    async handleReact(@MessageBody() data: ReactMessageOnCommand): Promise<void> {
        // this.server.to(data.chatID).emit('message', result);
    }

    @SubscribeMessage('reply')
    async handleReply(@MessageBody() data: ReplyOnMessageCommand): Promise<void> {
        // this.server.to(data.chatID).emit('message', result);
    }



    // @SubscribeMessage('userTyping')
    // async handleTyping(@MessageBody() data: any): Promise<void> {
    //     this.server.to(data.chatID).emit('message', result);
    // }
}