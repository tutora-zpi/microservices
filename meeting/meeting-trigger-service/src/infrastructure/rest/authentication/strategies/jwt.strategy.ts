import { ExtractJwt, Strategy } from 'passport-jwt';
import { PassportStrategy } from '@nestjs/passport';
import { ConfigService } from '@nestjs/config';
import { JWT_SECRET } from '../authentication.module';
import { NoJWTSecretProvided } from 'src/domain/exceptions/jwt-secret.exception';
import { Injectable } from '@nestjs/common';


@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
    constructor(configService: ConfigService) {
        const secret = configService.get<string>(JWT_SECRET);
        if (!secret) {
            throw new NoJWTSecretProvided();
        }

        super({
            jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
            secretOrKey: secret,
        });
    }

    async validate(payload: any) {
        return {
            // zakładam ze bedzie cos takiego
            userId: payload.sub,
            email: payload.email,
            role: payload.role
        };
    }
}
