import { ExtractJwt, Strategy } from 'passport-jwt';
import { PassportStrategy } from '@nestjs/passport';
import { JwksService } from './jwks.service';
import { Injectable, UnauthorizedException } from '@nestjs/common';
import * as jwt from 'jsonwebtoken';
import { IsUUID, IsEmail, IsString, IsNumber, IsOptional, validateOrReject } from 'class-validator';
import { plainToClass } from 'class-transformer';

export class Payload {
    @IsUUID()
    sub: string;

    @IsEmail()
    email: string;

    @IsString()
    roles: string;

    @IsNumber()
    @IsOptional()
    iat?: number;

    @IsNumber()
    exp: number;

    constructor(partial: Partial<Payload>) {
        Object.assign(this, partial);
    }
}

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

    async validate(payload: any) {
        const payloadInstance = plainToClass(Payload, payload);

        try {
            await validateOrReject(payloadInstance);
        } catch (errors) {
            throw new UnauthorizedException('Invalid token payload');
        }

        return {
            userId: payloadInstance.sub,
            email: payloadInstance.email,
            roles: payloadInstance.roles,
        };
    }
}