import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { MongooseModule } from '@nestjs/mongoose';

import { AppController } from './app.controller';
import { AppService } from './app.service';
import { AuthModule } from './auth/auth.module';
import { RoadmapModule } from './roadmap/roadmap.module';
import { TopicModule } from './topic/topic.module';

@Module({
  imports: [
    // Loads .env variables globally, such as JWT_SECRET and JWT_EXPIRES_IN_SECONDS.
    ConfigModule.forRoot({
      isGlobal: true,
    }),

    // Existing MongoDB connection used by the team backend.
    MongooseModule.forRoot('mongodb://localhost/roadmap2'),

    RoadmapModule,
    TopicModule,
    AuthModule,
  ],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {}
