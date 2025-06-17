import {
    WebSocketGateway,
    OnGatewayConnection,
    OnGatewayDisconnect,
    SubscribeMessage,
    MessageBody,
    ConnectedSocket,
} from '@nestjs/websockets';
import { Logger } from '@nestjs/common';
import { Socket } from 'socket.io';
import { AutosaveScheduler } from '../scheduler/autosave.scheduler'
import { WebSocketServer } from '@nestjs/websockets';
import { Server } from 'socket.io';

@WebSocketGateway({ cors: true })
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
    }

    @SubscribeMessage('join-session')
    handleJoinSession(
        @MessageBody() payload: { sessionId: string },
        @ConnectedSocket() client: Socket,
    ) {
        client.join(payload.sessionId);
        this.logger.log(`👥 ${client.id} joined session ${payload.sessionId}`);

        const room = this.server.sockets.adapter.rooms.get(payload.sessionId);
        this.logger.debug(`🔍 Room "${payload.sessionId}" has ${room?.size || 0} clients`);
    }

    @SubscribeMessage('leave-session')
    handleLeaveSession(
        @MessageBody() payload: { sessionId: string },
        @ConnectedSocket() client: Socket,
    ) {
        client.leave(payload.sessionId);
        this.logger.log(`👋 ${client.id} left session ${payload.sessionId}`);
    }

    @SubscribeMessage('board:update')
    handleBoardUpdate(
        @MessageBody() payload: { sessionId: string; data: any },
        @ConnectedSocket() client: Socket,
    ) {
        const { sessionId, data } = payload;
        this.logger.debug(`📥 board:update from ${client.id} (session: ${sessionId})`);

        this.autosave.bufferBoard(sessionId, data);

        client.to(sessionId).emit('board:sync', data);
        // this.server.to(sessionId).emit('board:sync', data);
    }
}
