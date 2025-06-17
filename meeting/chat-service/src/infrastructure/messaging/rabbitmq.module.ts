import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { RabbitMQSetupService } from './rabbitmq.setup.service';

@Module({
    imports: [ConfigModule],
    providers: [RabbitMQSetupService],
})
export class RabbitMQModule { }
