import { Module } from '@nestjs/common';
import { databaseProviders } from './database.provider';
import { repoProviders } from './repositories/repository.provider';

@Module({
    providers: [
        ...databaseProviders,
        ...repoProviders,
    ],
    exports: [
        ...databaseProviders,
        ...repoProviders,
    ],
})
export class DatabaseModule { }
