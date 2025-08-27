import { Injectable, NestInterceptor, ExecutionContext, CallHandler } from '@nestjs/common';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { ServiceResponse } from './response';
import { DTO } from '../dto/dto';

@Injectable()
export class ServiceResponseWrapperInterceptor<T extends DTO>
  implements NestInterceptor<T, ServiceResponse<T>> {

  intercept(context: ExecutionContext, next: CallHandler<T>): Observable<ServiceResponse<T>> {
    return next.handle().pipe(
      map(data => new ServiceResponse<T>(data)),
      catchError((err: any) => {
        const errorMessage = err?.response?.message || err?.message || 'Unknown error';
        return of(new ServiceResponse<T>(undefined, errorMessage));
      })
    );
  }
}
