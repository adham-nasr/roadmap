import { Module } from '@nestjs/common';
import { RoadmapService } from './roadmap.service';
import { RoadmapController } from './roadmap.controller';
import { TopicModule } from 'src/topic/topic.module';
import { MongooseModule } from '@nestjs/mongoose';
import { Roadmap, roadmapSchema } from './roadmap.schema';

@Module({
  imports: [
    MongooseModule.forFeature([{name:Roadmap.name,schema:roadmapSchema}]),
    TopicModule],
  controllers: [RoadmapController],
  providers: [RoadmapService],
})
export class RoadmapModule {}
