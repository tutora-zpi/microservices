import { ConfigService } from '@nestjs/config';
import { RmqOptions, Transport } from '@nestjs/microservices';

export const getRmqOptions = (configService: ConfigService): RmqOptions => ({
    transport: Transport.RMQ,
    options: {
        urls: [configService.get<string>('RABBITMQ_URL') || 'amqp://user:user@localhost:5672'],
        queue: 'meeting',
        exchange: 'meeting_events_exchange',
        exchangeType: 'fanout',
        routingKey: '',
        queueOptions: {
            durable: true,
        },
        noAck: false,
    },
});
