import { CanActivate, ExecutionContext, Injectable, Logger, UnauthorizedException } from '@nestjs/common';
import { Socket } from 'socket.io';
import * as jwt from 'jsonwebtoken';
import { JwtStrategy } from '../jwt.strategy';

@Injectable()
export class WsAuthGuard implements CanActivate {
  private readonly logger = new Logger(WsAuthGuard.name);
  private readonly algorithms: jwt.Algorithm[] = ['RS256'];

  constructor(private readonly jwtStrategy: JwtStrategy) { }

  async canActivate(context: ExecutionContext): Promise<boolean> {
    const client: Socket = context.switchToWs().getClient();
    const token = this.extractToken(client);

    if (!token) {
      this.logger.warn('Token not found');
      throw new UnauthorizedException('Token not provided');
    }

    try {
      const decoded: any = jwt.decode(token, { complete: true });
      if (!decoded?.header?.kid) {
        this.logger.warn('No kid in JWT header');
        throw new UnauthorizedException('Invalid token header');
      }

      const key = await this.jwtStrategy.jwksService.getSigningKey(decoded.header.kid);

      const payload = jwt.verify(token, key, { algorithms: this.algorithms }) as any;

      const validatedPayload = await this.jwtStrategy.validate(payload);

      client.data.user = validatedPayload;
      return true;
    } catch (err) {
      this.logger.error('Token verification failed', err);
      throw new UnauthorizedException('Invalid token');
    }
  }

  private extractToken(client: Socket): string | null {
    const header = client.handshake.headers.authorization;
    if (header && header.startsWith('Bearer ')) {
      return header.slice(7);
    }

    return (
      (client.handshake.query?.token as string | undefined) ||
      (client.handshake.auth?.token as string | undefined) ||
      null
    );
  }
}