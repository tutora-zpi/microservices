import { Board } from '../models/board.model';

export interface IBoardRepository {
    save(board: Board): Promise<void>;
    getLatest(sessionId: string): Promise<Board | null>;
}