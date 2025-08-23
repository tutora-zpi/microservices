import { Injectable, OnModuleInit, Logger } from '@nestjs/common';
import * as amqp from 'amqplib';
import { RabbitMQConfig } from '../config/rabbitmq.config';

@Injectable()
export class RabbitMQSetupService implements OnModuleInit {
  private readonly logger = new Logger(RabbitMQSetupService.name);

  constructor(
    private readonly rabbitmqConfig: RabbitMQConfig,
  ) { }

  async onModuleInit() {
    try {
      const url = this.rabbitmqConfig.url();

      this.logger.log(`Connecting with URL: ${url}`);
      const connection = await amqp.connect(url);

      const channel = await connection.createChannel();

      await channel.assertExchange(this.rabbitmqConfig.exchange, 'fanout', { durable: true });
      await channel.assertQueue(this.rabbitmqConfig.queue, { durable: true });
      await channel.bindQueue(this.rabbitmqConfig.queue, this.rabbitmqConfig.exchange, '');

      this.logger.log(`Queue "${this.rabbitmqConfig.queue}" bound to exchange "${this.rabbitmqConfig.exchange}"`);

      await channel.close();
      await connection.close();
    } catch (error) {
      this.logger.error('Failed to setup RabbitMQ bindings', error);
    }
  }
}
