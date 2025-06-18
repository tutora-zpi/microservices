import { Injectable, CanActivate, ExecutionContext, UnauthorizedException, Logger } from '@nestjs/common';
import { Request } from 'express';
import { JwtStrategy } from '../jwt.strategy';
import * as jwt from 'jsonwebtoken';

@Injectable()
export class AuthGuard implements CanActivate {
    private readonly userKey = 'user';
    private readonly logger = new Logger(AuthGuard.name);

    constructor(private readonly jwtStrategy: JwtStrategy) { }

    async canActivate(context: ExecutionContext): Promise<boolean> {
        const request = context.switchToHttp().getRequest<Request>();
        const authHeader = request.headers.authorization;

        if (!authHeader) {
            throw new UnauthorizedException('Authorization header missing');
        }

        const [bearer, token] = authHeader.split(' ');

        if (bearer !== 'Bearer' || !token) {
            throw new UnauthorizedException('Invalid authorization header format');
        }

        try {
            const decoded: any = jwt.decode(token, { complete: true });
            if (!decoded?.header?.kid) {
                throw new UnauthorizedException('Invalid token header');
            }

            const key = await this.jwtStrategy.jwksService.getSigningKey(decoded.header.kid);

            const payload = await new Promise((resolve, reject) => {
                jwt.verify(token, key, { algorithms: ['RS256'] }, (err, decoded) => {
                    if (err) {
                        reject(err);
                    } else {
                        resolve(decoded);
                    }
                });
            });

            const validatedPayload = await this.jwtStrategy.validate(payload);

            this.logger.debug("Saving user data under 'user'")
            request[this.userKey] = validatedPayload;

            return true;
        } catch (error) {
            console.log('error', error)
            throw new UnauthorizedException('Invalid or expired token');
        }
    }
}