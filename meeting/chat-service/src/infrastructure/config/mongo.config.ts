import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';

@Injectable()
export class MongoDBConfig {
    private readonly logger = new Logger(MongoDBConfig.name);

    constructor(private readonly configService: ConfigService) { }

    url(): string {
        let uri = this.configService.get<string>('MONGO_URI');

        if (!uri) {
            this.logger.error(
                'Missing MONGO_URI environment variable. Building URI from other variables...',
            );

            const user = this.configService.get<string>('MONGO_INITDB_ROOT_USERNAME');
            const password = this.configService.get<string>('MONGO_INITDB_ROOT_PASSWORD');
            const host = this.configService.get<string>('MONGO_HOST');
            const port = this.configService.get<string>('MONGO_PORT', '27017');
            const dbName = this.configService.get<string>('MONGO_DB_NAME', 'chat_db');

            if (!user || !password || !host || !port || !dbName) {
                this.logger.error(
                    'Missing one or more environment variables for MongoDB connection: ' +
                    'MONGO_INITDB_ROOT_USERNAME, MONGO_INITDB_ROOT_PASSWORD, MONGO_HOST, MONGO_PORT, MONGO_DB_NAME',
                );

                throw new URIError(
                    'Failed to connect with database, please provide valid URI.',
                );
            }

            uri = `mongodb://${user}:${password}@${host}:${port}/${dbName}?authSource=admin`;
            this.logger.log(`Constructed MONGO_URI: ${uri}`);
        }

        return uri;
    }
}
