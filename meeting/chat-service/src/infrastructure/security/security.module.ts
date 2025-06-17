import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { JwtStrategy } from './jwt.strategy';
import { JwksService } from './jwks.service';
import { WsAuthGuard } from './guards/ws.auth.guard';
import { AuthGuard } from './guards/auth.guard';

@Module({
    imports: [ConfigModule],
    providers: [JwtStrategy, JwksService, WsAuthGuard, AuthGuard],
    exports: [WsAuthGuard, AuthGuard, JwtStrategy, JwksService],
})
export class SecurityModule { }