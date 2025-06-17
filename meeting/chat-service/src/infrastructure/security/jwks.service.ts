import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import * as jwksClient from 'jwks-rsa';

@Injectable()
export class JwksService {
    private readonly logger: Logger = new Logger(JwksService.name)

    private client: jwksClient.JwksClient;

    constructor(private configService: ConfigService) {
        this.client = jwksClient({
            jwksUri: this.configService.get<string>('JWKS_URL') || "http://localhost:8080/.well-known/jwks.json",
            cache: true,
            cacheMaxEntries: 5,
            cacheMaxAge: 3600000,
            rateLimit: true,
            jwksRequestsPerMinute: 10,
        });
    }

    async getSigningKey(kid: string): Promise<string> {
        this.logger.debug(`Fetching signing key for kid: ${kid}`);
        return new Promise((resolve, reject) => {
            this.client.getSigningKey(kid, (err, key) => {
                if (err) {
                    this.logger.error(`Error fetching signing key: ${err.message}`, err.stack);
                    return reject(err);
                }
                if (!key) {
                    this.logger.error('No key returned from JWKS endpoint');
                    return reject(new Error('No key found'));
                }
                const signingKey = key.getPublicKey();
                this.logger.debug(`Obtained signing key: ${signingKey}`);
                resolve(signingKey);
            });
        });
    }

}
