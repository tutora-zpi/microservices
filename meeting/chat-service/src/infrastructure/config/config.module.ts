import { Module } from '@nestjs/common';
import { ConfigModule as CM } from '@nestjs/config';
import { MongoDBConfig } from './mongo.config';
import { RabbitMQConfig } from './rabbitmq.config';


@Module({
    imports: [CM],
    providers: [MongoDBConfig, RabbitMQConfig],
    exports: [MongoDBConfig, RabbitMQConfig],
})
export class ConfigModule { }
