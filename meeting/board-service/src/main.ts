import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { RabbitMQConfig } from './infrastructure/messaging/rabbitmq.config';


async function bootstrap() {
    const app = await NestFactory.create(AppModule, { cors: true });

    const rabbitConfig = app.get(RabbitMQConfig);

    app.connectMicroservice(rabbitConfig.options());

    await app.startAllMicroservices();
    await app.listen(process.env.PORT ?? 8001);

    console.log('Listening on port', process.env.PORT ?? 8001);
}

bootstrap();