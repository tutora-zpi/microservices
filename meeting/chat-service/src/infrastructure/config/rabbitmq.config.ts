import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { RmqOptions, Transport } from '@nestjs/microservices';
import { URINotFound } from 'src/domain/exceptions/uri.exception';

@Injectable()
export class RabbitMQConfig {
    private readonly logger = new Logger(RabbitMQConfig.name);

    readonly exchange: string;
    readonly queue: string;

    constructor(private readonly configService: ConfigService) {
        this.exchange =
            this.configService.get<string>('EVENT_EXCHANGE_QUEUE_NAME', 'meeting_events_exchange');

        this.queue =
            this.configService.get<string>('QUEUE_NAME', 'meeting');
    }

    url(): string {
        const url = this.configService.get<string>('RABBITMQ_URL');
        if (url) {
            this.logger.log('Using RABBITMQ_URL from environment');
            return url;
        }

        const host = this.configService.get<string>('RABBITMQ_HOST');
        const user = this.configService.get<string>('RABBITMQ_DEFAULT_USER');
        const pass = this.configService.get<string>('RABBITMQ_DEFAULT_PASS');
        const port = this.configService.get<string>('RABBITMQ_PORT');

        if (!host || !user || !pass || !port) {
            this.logger.error(
                'Missing one or more RabbitMQ environment variables: RABBITMQ_HOST, RABBITMQ_DEFAULT_USER, RABBITMQ_DEFAULT_PASS, RABBITMQ_PORT',
            );
            throw new URINotFound('Please fill .env.local file, especially RabbitMQ section!');
        }

        const constructedUrl = `amqp://${user}:${pass}@${host}:${port}/`;
        this.logger.log(`Constructed RabbitMQ URL: ${constructedUrl}`);
        return constructedUrl;
    }


    options(): RmqOptions {
        return {
            transport: Transport.RMQ,
            options: {
                urls: [this.url()],
                queue: this.queue,
                exchange: this.exchange,
                exchangeType: 'fanout',
                routingKey: '',
                queueOptions: {
                    durable: true,
                },
                noAck: false,
            },
        };
    }
}
