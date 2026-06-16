import { Module } from '@nestjs/common';
import { RoadmapService } from './roadmap.service';
import { RoadmapController } from './roadmap.controller';
import { TopicModule } from 'src/topic/topic.module';
import { MongooseModule } from '@nestjs/mongoose';
import { Roadmap, roadmapSchema } from './roadmap.schema';
import { RoadmapRepository } from './roadmap.repository';

@Module({
  imports: [
    MongooseModule.forFeature([{ name: Roadmap.name, schema: roadmapSchema }]),
    TopicModule,
  ],
  controllers: [RoadmapController],
  providers: [RoadmapService, RoadmapRepository],
})
export class RoadmapModule {}
