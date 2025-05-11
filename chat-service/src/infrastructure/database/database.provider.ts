import { Logger, Provider } from '@nestjs/common';
import * as mongoose from 'mongoose';

export const DATABASE_CONNECTION = 'DATABASE_CONNECTION';

const logger = new Logger('DatabaseProvider');

export const databaseProviders: Provider[] = [
    {
        provide: DATABASE_CONNECTION,
        useFactory: async (): Promise<typeof mongoose> => {
            const uri = process.env.MONGO_URI;

            if (!uri) {
                logger.error('Missing MONGO_URI environment variable.');
                throw new Error('Failed to connect with database, please provide valid URI.');
            }

            try {
                const connection = await mongoose.connect(uri);
                logger.log('Successfully connected to MongoDB.');
                return connection;
            } catch (err) {
                logger.error('MongoDB connection failed', err);
                throw err;
            }
        },
    },
];
