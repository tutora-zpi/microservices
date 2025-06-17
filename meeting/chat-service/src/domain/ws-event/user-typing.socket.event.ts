import { SocketEvent } from "./socket.event";

export class UserTyping implements SocketEvent {
    constructor(
        public readonly chatID: string,
        public readonly userID: string,
        public readonly isTyping: boolean
    ) { }
}