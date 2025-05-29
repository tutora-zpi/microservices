export class NoJWTSecretProvided extends Error {
    constructor(message: string = 'JWT_SECRET is not defined in the configuration. Please provide it in .env* file/s') {
        super(message);
        this.name = 'NoJWTSecretProvided';
    }
}