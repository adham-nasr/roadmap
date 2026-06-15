import { Module } from '@nestjs/common';
import { TopicService } from './topic.service';
import { TopicController } from './topic.controller';
import { MongooseModule } from '@nestjs/mongoose';
import { Topic, topicSchema } from './topic.schema';
import { TopicRepository } from './topic.repository';

@Module({
  imports:[MongooseModule.forFeature([{name:Topic.name,schema:topicSchema}]) 
  ],
  controllers: [TopicController],
  providers: [TopicService,TopicRepository],
  exports: [MongooseModule,TopicRepository]
})
export class TopicModule {}
