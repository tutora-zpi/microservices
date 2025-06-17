import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';

import { AppConfigModule } from './config/config.module';
import { DatabaseModule } from './infrastructure/database/database.module';
import { BoardService } from './app/services/board.service';
import { BoardRepositoryImpl } from './infrastructure/persistence/board.repository.impl';
import { BoardSchema, BoardEntity } from './infrastructure/database/schemas/board.schema';
import { BoardController } from './app/controllers/board.controller';
import { ScheduleModule } from '@nestjs/schedule';
import { AutosaveScheduler } from './infrastructure/scheduler/autosave.scheduler';
import { BoardGateway } from './infrastructure/ws/board.gateway'


@Module({
    imports: [
        ScheduleModule.forRoot(),
        AppConfigModule,
        DatabaseModule,
        MongooseModule.forFeature([{ name: BoardEntity.name, schema: BoardSchema }])
    ],
    providers: [
        AutosaveScheduler,
        BoardGateway,
        BoardService,
        {
            provide: 'BoardRepository',
            useClass: BoardRepositoryImpl,
        }
    ],
    controllers: [BoardController],
    exports: [BoardService],
})
export class AppModule {}
