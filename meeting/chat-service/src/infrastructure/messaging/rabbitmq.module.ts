import { Module } from '@nestjs/common';
import { ClientsModule } from '@nestjs/microservices';
import { ConfigModule } from '../config/config.module';
import { RabbitMQConfig } from '../config/rabbitmq.config';


@Module({
  imports: [
    ConfigModule,
    ClientsModule.registerAsync([
      {
        imports: [ConfigModule],
        name: 'RABBITMQ_SERVICE',
        useFactory: (config: RabbitMQConfig) => config.options(),
        inject: [RabbitMQConfig],
      },
    ]),
  ],
  providers: [ClientsModule],
  exports: [ClientsModule],
})

export class RabbitMQModule { }