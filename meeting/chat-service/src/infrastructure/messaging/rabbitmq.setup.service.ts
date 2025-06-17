import { Injectable, OnModuleInit, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import * as amqp from 'amqplib';

// Used to bind queue
@Injectable()
export class RabbitMQSetupService implements OnModuleInit {
    private readonly logger = new Logger(RabbitMQSetupService.name);

    constructor(private readonly configService: ConfigService) { }

    async onModuleInit() {
        const url = this.configService.get<string>('RABBITMQ_URL') || 'amqp://user:user@localhost:5672';
        const exchange = this.configService.get<string>('EVENT_EXCHANGE_QUEUE_NAME') || 'meeting_events_exchange';
        const queue = this.configService.get<string>('QUEUE_NAME') || 'meeting';

        try {
            const connection = await amqp.connect(url);
            const channel = await connection.createChannel();

            await channel.assertExchange(exchange, 'fanout', { durable: true });
            await channel.assertQueue(queue, { durable: true });
            await channel.bindQueue(queue, exchange, '');

            this.logger.log(`Queue "${queue}" bound to exchange "${exchange}"`);

            await channel.close();
            await connection.close();
        } catch (error) {
            this.logger.error('Failed to setup RabbitMQ bindings', error);
        }
    }
}
