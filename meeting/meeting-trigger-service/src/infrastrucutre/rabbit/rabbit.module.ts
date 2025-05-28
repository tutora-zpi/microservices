import { ConfigModule, ConfigService } from "@nestjs/config";
import { ClientsModule } from "@nestjs/microservices";
import { Module } from "@nestjs/common";
import { getRmqOptions } from "src/config /rabbit.config";

export const RABBITMQ_SERVICE = 'RABBITMQ_SERVICE';

@Module({
    imports: [
        ConfigModule,
        ClientsModule.registerAsync([
            {
                name: 'RABBITMQ_SERVICE',
                imports: [ConfigModule],
                useFactory: (configService: ConfigService) => getRmqOptions(configService),
                inject: [ConfigService],
            },
        ]),
    ],
    exports: [ClientsModule],
})

export class RabbitModule { }
