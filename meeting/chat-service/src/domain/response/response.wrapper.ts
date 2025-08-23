import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from '@nestjs/common';
import { map } from 'rxjs';
import { ServiceResponse } from './response';
import { DTO } from '../dto/dto';

@Injectable()
export class ServiceResponseWrapperInterceptor<T extends DTO>
  implements NestInterceptor<T, ServiceResponse<T>> {
  intercept(context: ExecutionContext, next: CallHandler) {
    return next.handle().pipe(
      map((data) => new ServiceResponse(data)),
    );
  }
}
