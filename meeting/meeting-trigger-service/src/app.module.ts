import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { ClientsModule } from '@nestjs/microservices';
import { getRmqOptions } from './config /rabbit.config';

export const RABBITMQ_SERVICE = 'RABBITMQ_SERVICE';

@Module({
  imports: [
    ClientsModule.registerAsync([
      {
        name: RABBITMQ_SERVICE,
        imports: [ConfigModule],
        useFactory: (configService: ConfigService) => getRmqOptions(configService),
        inject: [ConfigService],
      },
    ]),

  ],
  controllers: [],
  providers: [],
})
export class AppModule { }
