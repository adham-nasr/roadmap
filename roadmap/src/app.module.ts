import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { AppService } from './app.service';
import { RoadmapModule } from './roadmap/roadmap.module';
import { TopicModule } from './topic/topic.module';
import { ResourceModule } from './resource/resource.module';
import { MongooseModule } from '@nestjs/mongoose';

@Module({
  imports: [ MongooseModule.forRoot('mongodb://localhost/roadmap'),

  RoadmapModule, TopicModule, ResourceModule],
  controllers: [AppController],
  providers: [AppService],
  exports:[MongooseModule]
})
export class AppModule {}
