import { Logger } from '@nestjs/common';
import { CommandBus } from '@nestjs/cqrs';
import { WebSocketGateway, WebSocketServer, OnGatewayConnection, OnGatewayDisconnect, MessageBody, SubscribeMessage } from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';
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

    @SubscribeMessage('sendMessage')
    async handleMessage(@MessageBody() data: SendMessageCommand): Promise<void> {
        this.server.emit('message', data);

        const d = new SendMessageCommand(data.content, data.senderID, data.receiverID);

        const result = await this.commandBus.execute(d);

        this.server.emit(result);
    }
}