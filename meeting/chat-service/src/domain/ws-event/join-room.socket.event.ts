import { SocketEvent } from "./socket.event";

export class JoinToRoomEvent implements SocketEvent {
    constructor(
        public readonly roomId: string,
        public readonly token: string,
    ) { }
}