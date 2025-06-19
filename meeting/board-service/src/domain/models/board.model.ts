export class Board {
    constructor(
        public sessionId: string,
        public excalidrawData: any,
        public updatedAt: Date = new Date(),
    ) {}
}