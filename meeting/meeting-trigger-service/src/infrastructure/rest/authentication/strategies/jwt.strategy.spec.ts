import { ConfigService } from '@nestjs/config';
import { JwtStrategy } from './jwt.strategy';
import { NoJWTSecretProvided } from 'src/domain/exceptions/jwt-secret.exception';


describe('JwtStrategy', () => {
    let configService: ConfigService;

    beforeEach(() => {
        configService = {
            get: jest.fn()
        } as any as ConfigService;
    });

    it('should throw NoJWTSecretProvided if secret is not provided', () => {
        (configService.get as jest.Mock).mockReturnValue(undefined);

        expect(() => new JwtStrategy(configService)).toThrow(NoJWTSecretProvided);
    });

    it('should create an instance if secret is provided', () => {
        (configService.get as jest.Mock).mockReturnValue('supersecret');

        const strategy = new JwtStrategy(configService);

        expect(strategy).toBeDefined();
    });

    it('validate() should return user data from payload', async () => {
        (configService.get as jest.Mock).mockReturnValue('supersecret');
        const strategy = new JwtStrategy(configService);

        const payload = { sub: '123', email: 'test@example.com', role: 'user' };
        const result = await strategy.validate(payload);

        expect(result).toEqual({
            userId: '123',
            email: 'test@example.com',
            role: 'user',
        });
    });
});
