import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { BoardService } from 'src/app/services/board.service';
import { AutosaveScheduler } from 'src/infrastructure/scheduler/autosave.scheduler';
import { BoardUpdateHandler } from './board-update-handler';
import { BoardRepositoryImpl } from 'src/infrastructure/persistence/board.repository.impl';
import { BoardEntity, BoardSchema } from 'src/infrastructure/database/schemas/board.schema';

@Module({
    imports: [
        MongooseModule.forFeature([{ name: BoardEntity.name, schema: BoardSchema }]),
    ],
    controllers: [BoardUpdateHandler],
    providers: [
        BoardService,
        AutosaveScheduler,
        {
            provide: 'Scheduler',
            useExisting: AutosaveScheduler,
        },
        {
            provide: 'BoardRepository',
            useClass: BoardRepositoryImpl,
        },
    ],
    exports: ['Scheduler', 'BoardRepository'],
})
export class BoardModule { }

