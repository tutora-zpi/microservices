import { Module } from '@nestjs/common';
import { databaseProviders } from './database.provider';
import { repoProviders } from './repositories/repository.provider';
import { MongoDBConfig } from '../config/mongo.config';

@Module({
  providers: [...databaseProviders, ...repoProviders, MongoDBConfig],
  exports: [...databaseProviders, ...repoProviders],
})
export class DatabaseModule { }
