import { Injectable, Logger } from '@nestjs/common';
import { Interval } from '@nestjs/schedule';
import { BoardService } from 'src/app/services/board.service';
import { Scheduler } from './scheduler.interface';

@Injectable()
export class AutosaveScheduler implements Scheduler {
    private readonly logger = new Logger(AutosaveScheduler.name);
    private boardBuffer = new Map<string, any>();

    constructor(private readonly boardService: BoardService) { }

    bufferBoard(sessionId: string, data: any) {
        this.boardBuffer.set(sessionId, data);
    }

    @Interval(10000)
    async handleAutosave() {
        if (this.boardBuffer.size === 0) return;

        this.logger.log(`Autosaving ${this.boardBuffer.size} sessions...`);

        for (const [sessionId, data] of this.boardBuffer.entries()) {
            try {
                await this.boardService.saveBoard(sessionId, data);
                this.logger.log(`Saved board for session: ${sessionId}`);
            } catch (err) {
                this.logger.error(`Failed to save board ${sessionId}`, err);
            }
        }

        this.boardBuffer.clear();
    }

    getBuffer(sessionId: string): any {
        return this.boardBuffer.get(sessionId);
    }

    async flushSingle(sessionId: string): Promise<void> {
        const data = this.boardBuffer.get(sessionId);
        if (!data) return;

        await this.boardService.saveBoard(sessionId, data);
        this.boardBuffer.delete(sessionId);
    }

    async getBoard(sessionId: string) {
        return await this.boardService.getBoard(sessionId);
    }

    async saveNow(sessionId: string, data: any) {
        try {
            await this.boardService.saveBoard(sessionId, data);
            this.logger.log(`Immediate save for session ${sessionId}`);
            this.boardBuffer.delete(sessionId);
        } catch (err) {
            this.logger.error(`Immediate save failed for session ${sessionId}`, err);
        }
    }
}
