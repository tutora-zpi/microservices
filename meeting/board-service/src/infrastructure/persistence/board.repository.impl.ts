import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Model } from 'mongoose';

import { BoardEntity, BoardDocument } from '../database/schemas/board.schema';
import { IBoardRepository } from '../../domain/repositories/board.repository';
import { Board } from '../../domain/models/board.model';

@Injectable()
export class BoardRepositoryImpl implements IBoardRepository {
    constructor(
        @InjectModel(BoardEntity.name) private readonly boardModel: Model<BoardDocument>
    ) {}

    async save(board: Board): Promise<void> {
        await this.boardModel.findOneAndUpdate(
            { sessionId: board.sessionId },
            {
                sessionId: board.sessionId,
                excalidrawData: board.excalidrawData,
                updatedAt: new Date(),
            },
            { upsert: true }
        );
    }

    async getLatest(sessionId: string): Promise<Board | null> {
        const doc = await this.boardModel.findOne({ sessionId }).sort({ updatedAt: -1 });
        if (!doc) return null;
        return new Board(doc.sessionId, doc.excalidrawData, doc.updatedAt);
    }
}