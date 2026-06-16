import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { RoadmapModule } from './roadmap/roadmap.module';
import { TopicModule } from './topic/topic.module';
import { MongooseModule } from '@nestjs/mongoose';

@Module({
  imports: [
    MongooseModule.forRoot(
      'mongodb+srv://ITI-FinalProject:ITI-PROJECT-123@cluster0.1rl7if2.mongodb.net/roadmapsdb',
    ),
    RoadmapModule,
    TopicModule,
  ],
  controllers: [AppController],
  providers: [AppService],
  exports: [MongooseModule],
})
export class AppModule {}
