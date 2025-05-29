import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { ConfigService } from '@nestjs/config';
import { BoardEntity, BoardSchema } from './schemas/board.schema';

@Module({
    imports: [
        MongooseModule.forRootAsync({
            useFactory: (configService: ConfigService) => ({
                uri: configService.get<string>('MONGO_URI') || 'mongodb://localhost:27017/board',
            }),
            inject: [ConfigService],
        }),
            MongooseModule.forFeature([
            { name: BoardEntity.name, schema: BoardSchema },
        ]),
    ],
    exports: [MongooseModule], 
})
export class DatabaseModule {}