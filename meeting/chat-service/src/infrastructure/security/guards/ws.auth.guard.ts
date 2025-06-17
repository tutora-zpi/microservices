import { CanActivate, ExecutionContext, Injectable, UnauthorizedException } from '@nestjs/common';
import { Socket } from 'socket.io';
import * as jwt from 'jsonwebtoken';
import { JwtStrategy } from '../jwt.strategy';

@Injectable()
export class WsAuthGuard implements CanActivate {
    constructor(private readonly jwtStrategy: JwtStrategy) { }

    async canActivate(context: ExecutionContext): Promise<boolean> {
        const client: Socket = context.switchToWs().getClient();

        const token = this.extractToken(client);
        if (!token) {
            throw new UnauthorizedException('Token not provided');
        }

        try {
            const decoded: any = jwt.decode(token, { complete: true });
            if (!decoded?.header?.kid) {
                throw new UnauthorizedException('Invalid token header');
            }

            const key = await this.jwtStrategy.jwksService.getSigningKey(decoded.header.kid);

            return new Promise((resolve, reject) => {
                jwt.verify(token, key, { algorithms: ['RS256'] }, async (err, payload) => {
                    if (err) {
                        return reject(new UnauthorizedException('Token verification failed'));
                    }

                    try {
                        const validatedPayload = await this.jwtStrategy.validate(payload);
                        client.data.user = validatedPayload;
                        resolve(true);
                    } catch (error) {
                        reject(new UnauthorizedException('Invalid token payload'));
                    }
                });
            });
        } catch (error) {
            throw new UnauthorizedException('Token processing failed');
        }
    }

    private extractToken(client: Socket): string | null {
        const authHeader = client.handshake.headers.authorization as string | undefined;
        if (authHeader?.startsWith('Bearer ')) {
            return authHeader.slice(7);
        }

        const tokenFromQuery = client.handshake.query?.token as string | undefined;
        if (tokenFromQuery) {
            return tokenFromQuery;
        }

        const authFromQuery = client.handshake.auth?.token as string | undefined;
        return authFromQuery ?? null;
    }
}