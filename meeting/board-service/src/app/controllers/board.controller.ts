import { Controller, Get, Param, Post, Body, HttpException, HttpStatus } from '@nestjs/common';
import { BoardService } from '../services/board.service';

@Controller('board')
export class BoardController {
    constructor(private readonly boardService: BoardService) {}

    @Get(':sessionId')
    async getBoard(@Param('sessionId') sessionId: string) {
        const board = await this.boardService.getBoard(sessionId);
        if (!board) {
            throw new HttpException('Board not found', HttpStatus.NOT_FOUND);
        }
        return board;
    }

    @Post(':sessionId')
    async saveBoard(@Param('sessionId') sessionId: string, @Body() body: any) {
        if (!body || Object.keys(body).length === 0) {
            throw new HttpException('Invalid board data', HttpStatus.BAD_REQUEST);
        }
        await this.boardService.saveBoard(sessionId, body);
        return { status: 'ok' };
    }
}