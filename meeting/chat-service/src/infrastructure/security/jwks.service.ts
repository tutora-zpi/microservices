import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import * as jwksClient from 'jwks-rsa';

@Injectable()
export class JwksService {
  private client: jwksClient.JwksClient;

  constructor(configService: ConfigService) {
    const jwksUri =
      configService.get<string>('JWKS_URL') ??
      'http://localhost:8080/.well-known/jwks.json';

    this.client = jwksClient({
      jwksUri,
      cache: true,
      cacheMaxEntries: 5,
      cacheMaxAge: 3600000,
      rateLimit: true,
      jwksRequestsPerMinute: 10,
    });
  }

  async getSigningKey(kid: string): Promise<string> {
    const key = await this.client.getSigningKey(kid);
    return key.getPublicKey();
  }
}
