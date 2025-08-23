import { Logger, Provider } from '@nestjs/common';
import * as mongoose from 'mongoose';
import { URINotFound } from 'src/domain/exceptions/uri.exception';
import { MongoDBConfig } from '../config/mongo.config';

export const DATABASE_CONNECTION = 'DATABASE_CONNECTION';
const serverSelectionTimeoutMS = 5000;

const logger = new Logger('DatabaseProvider');

export const databaseProviders: Provider[] = [
  {
    provide: DATABASE_CONNECTION,
    inject: [MongoDBConfig],
    useFactory: async (mongoDBURL: MongoDBConfig): Promise<typeof mongoose> => {
      try {
        const uri = mongoDBURL.url();

        const connection = await mongoose.connect(uri, {
          serverSelectionTimeoutMS,
        });

        logger.log('Successfully connected to MongoDB.');
        return connection;
      } catch (err) {
        logger.error('MongoDB connection failed', err);
        throw new mongoose.MongooseError('Failed to connect with mongo db');
      }
    },
  },
];
