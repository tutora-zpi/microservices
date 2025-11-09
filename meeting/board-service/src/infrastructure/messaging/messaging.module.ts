import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { ClientsModule } from '@nestjs/microservices';
import { RabbitMQConfig } from './rabbitmq.config';

@Module({
    imports: [
        ConfigModule,
        ClientsModule.registerAsync([
            {
                imports: [ConfigModule, RabbitMQModule],
                inject: [RabbitMQConfig],
                useFactory: (config: RabbitMQConfig) => config.options(),
                name: 'RABBITMQ_SERVICE',
            },
        ]),
    ],
    providers: [RabbitMQConfig],
    exports: [ClientsModule, RabbitMQConfig],
})
export class RabbitMQModule { }
