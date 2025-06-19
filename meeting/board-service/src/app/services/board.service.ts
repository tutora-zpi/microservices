import { Injectable, Inject } from '@nestjs/common';
import { IBoardRepository } from '../../domain/repositories/board.repository';
import { Board } from '../../domain/models/board.model';

@Injectable()
export class BoardService {
    constructor(
        @Inject('BoardRepository') private readonly repo: IBoardRepository
    ) {}

    async saveBoard(sessionId: string, data: any): Promise<void> {
        const board = new Board(sessionId, data);
        await this.repo.save(board);
    }

    async getBoard(sessionId: string): Promise<Board | null> {
        return this.repo.getLatest(sessionId);
    }
}