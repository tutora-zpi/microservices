import { NestFactory } from '@nestjs/core';
import { SwaggerModule, DocumentBuilder } from '@nestjs/swagger';
import { AppModule } from './app.module';
import {
  ConsoleLogger,
  INestApplication,
  ValidationPipe,
} from '@nestjs/common';
import { IoAdapter } from '@nestjs/platform-socket.io';
import { MicroserviceOptions } from '@nestjs/microservices';
import { ConfigService } from '@nestjs/config';
import { RabbitMQConfig } from './infrastructure/config/rabbitmq.config';
import { ServiceResponseWrapperInterceptor } from './domain/response/response.wrapper';
import { AllExceptionsFilter } from './infrastructure/filters/all-exception.filter';

const appName = `${process.env.APP_NAME || 'PROVIDE APP NAME'} - Chat Service`;

function swag(app: INestApplication) {
  const config = new DocumentBuilder()
    .setTitle(appName)
    .setDescription(
      'Description for service with all endpoints and sockets evnets.',
    )
    .setVersion('1.0')
    .build();

  const documentFactory = () => SwaggerModule.createDocument(app, config);
  SwaggerModule.setup('api/v1/docs', app, documentFactory);
}

async function bootstrap() {
  const app = await NestFactory.create(AppModule, {
    logger: new ConsoleLogger({
      prefix: appName,
    }),
    cors: true,
  });

  swag(app);
  app.useWebSocketAdapter(new IoAdapter(app));

  const rabbitmqConfig = new RabbitMQConfig(app.get(ConfigService));

  app.connectMicroservice<MicroserviceOptions>(
    rabbitmqConfig.options(),
  );

  app.useGlobalFilters(new AllExceptionsFilter());

  app.useGlobalPipes(
    new ValidationPipe({
      whitelist: true,
      forbidNonWhitelisted: true,
      transform: true,
    }),
  );

  await app.startAllMicroservices();

  app.useGlobalInterceptors(new ServiceResponseWrapperInterceptor());

  await app.listen(process.env.PORT ?? 8002);
}
bootstrap();
