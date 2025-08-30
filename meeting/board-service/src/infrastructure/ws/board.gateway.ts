import {
    WebSocketGateway,
    OnGatewayConnection,
    OnGatewayDisconnect,
    SubscribeMessage,
    MessageBody,
    ConnectedSocket,
    WebSocketServer,
} from '@nestjs/websockets';
import { Logger } from '@nestjs/common';
import { Socket, Server } from 'socket.io';
import { AutosaveScheduler } from '../scheduler/autosave.scheduler';

@WebSocketGateway({
    cors: {
        origin: '*',
    },
})
export class BoardGateway implements OnGatewayConnection, OnGatewayDisconnect {
    @WebSocketServer()
    server: Server;

    private readonly logger = new Logger(BoardGateway.name);

    constructor(private readonly autosave: AutosaveScheduler) {}

    handleConnection(client: Socket) {
        this.logger.log(`🟢 Client connected: ${client.id}`);
    }

    handleDisconnect(client: Socket) {
        this.logger.log(`🔴 Client disconnected: ${client.id}`);

        this.server.sockets.adapter.rooms.forEach((clients, sessionId) => {
            if (clients.has(client.id)) {
                const data = this.autosave.getBuffer(sessionId);
                if (data) {
                    this.autosave.saveNow(sessionId, data);
                }
            }
        });
    }

    @SubscribeMessage('join-session')
    async handleJoinSession(
        @MessageBody() payload: { sessionId: string },
        @ConnectedSocket() client: Socket,
    ) {
        const { sessionId } = payload;

        await client.join(sessionId);
        this.logger.log(`👥 ${client.id} joined session ${sessionId}`);

        const room = this.server.sockets.adapter.rooms.get(sessionId);
        this.logger.debug(
            `🔍 Room "${sessionId}" has ${room?.size || 0} clients`,
        );

        const bufferedData = this.autosave.getBuffer(sessionId);
        if (bufferedData) {
            await this.autosave.flushSingle(sessionId);
            this.logger.log(`💾 Flushed buffer to DB for ${sessionId}`);
        }

        const board = await this.autosave.getBoard(sessionId);
        if (board) {
            client.emit('board:sync', board.excalidrawData);
            this.logger.debug(`📤 Sent board data to ${client.id}`);
        } else {
            this.logger.debug(`⚠️ No board found in DB for ${sessionId}`);
        }
    }

    @SubscribeMessage('leave-session')
    async handleLeaveSession(
        @MessageBody() payload: { sessionId: string },
        @ConnectedSocket() client: Socket,
    ) {
        client.leave(payload.sessionId);
        this.logger.log(`👋 ${client.id} left session ${payload.sessionId}`);

        const data = this.autosave.getBuffer(payload.sessionId);
        if (data) {
            await this.autosave.saveNow(payload.sessionId, data);
        }
    }

    @SubscribeMessage('board:update')
    handleBoardUpdate(
        @MessageBody() payload: { sessionId: string; data: any },
        @ConnectedSocket() client: Socket,
    ) {
        const { sessionId, data } = payload;
        this.logger.debug(
            `📥 board:update from ${client.id} (session: ${sessionId})`,
        );

        this.autosave.bufferBoard(sessionId, data);

        client.to(sessionId).emit('board:sync', data);
        // this.server.to(sessionId).emit('board:sync', data);
    }
}
