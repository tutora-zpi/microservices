export interface Scheduler {

    bufferBoard(sessionId: string, data: any): void;
    getBuffer(sessionId: string): any;
    flushSingle(sessionId: string): Promise<void>;
    getBoard(sessionId: string): Promise<any>;
    saveNow(sessionId: string, data: any): Promise<void>;
}
