import { ConfigService } from '@nestjs/config';
import { RmqOptions, Transport } from '@nestjs/microservices';

const queueName = 'meeting';

export const getRmqOptions = (configService: ConfigService): RmqOptions => ({
    transport: Transport.RMQ,
    options: {
        urls: [configService.get<string>('RABBITMQ_URL') || 'amqp://user:user@localhost:5672'],
        queue: queueName,
        queueOptions: { durable: true },
        noAck: true,
    },
});
