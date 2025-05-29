import { Module } from '@nestjs/common';
import { AppConfigModule } from './config/config.module';
import { DatabaseModule } from './infrastructure/database/database.module';

@Module({
    imports: [
        AppConfigModule,
        DatabaseModule,
    ],
})
export class AppModule {}