import { ExceptionFilter, Catch, ArgumentsHost, HttpException, HttpStatus } from '@nestjs/common';
import { Response } from 'express';
import { ServiceResponse } from 'src/domain/response/response';

@Catch()
export class AllExceptionsFilter implements ExceptionFilter {
    catch(exception: any, host: ArgumentsHost) {
        const ctx = host.switchToHttp();
        const response = ctx.getResponse<Response>();
        const status =
            exception instanceof HttpException ? exception.getStatus() : HttpStatus.INTERNAL_SERVER_ERROR;
        const message =
            exception instanceof HttpException
                ? exception.message || (exception.getResponse() as any)?.message
                : 'Internal server error';

        response.status(status).json(new ServiceResponse(undefined, message));
    }
}
