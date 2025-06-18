import { ExtractJwt, Strategy } from 'passport-jwt';
import { PassportStrategy } from '@nestjs/passport';
import { JwksService } from './jwks.service';
import { Injectable, UnauthorizedException } from '@nestjs/common';
import * as jwt from 'jsonwebtoken';

export type PayloadUser = {
    userID: string;
    email: string;
    role: string;
};


@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
    constructor(readonly jwksService: JwksService) {
        super({
            jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
            secretOrKeyProvider: async (_, rawJwtToken, done) => {
                try {
                    const decoded: any = jwt.decode(rawJwtToken, { complete: true });
                    if (!decoded || !decoded.header?.kid) {
                        return done(new Error('No kid found in token'), null);
                    }
                    const key = await this.jwksService.getSigningKey(decoded.header.kid);
                    done(null, key);
                } catch (err) {
                    done(err, null);
                }
            },
            algorithms: ['RS256'],
        });
    }

    async validate(payload: any): Promise<PayloadUser> {
        if (!payload?.sub || !payload?.email || !payload?.roles) {
            throw new UnauthorizedException('Invalid JWT payload');
        }

        return {
            userID: payload.sub,
            email: payload.email,
            role: payload.roles
        }
    }
}