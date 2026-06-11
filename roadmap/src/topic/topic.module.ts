import { Module } from '@nestjs/common';
import { TopicService } from './topic.service';
import { TopicController } from './topic.controller';
import { MongooseModule } from '@nestjs/mongoose';
import { Topic, topicSchema } from './topic.schema';
import { ResourceModule } from 'src/resource/resource.module';

@Module({
  imports:[MongooseModule.forFeature([{name:Topic.name,schema:topicSchema}]) ,
          ResourceModule
  ],
  controllers: [TopicController],
  providers: [TopicService],
  exports: [MongooseModule]
})
export class TopicModule {}
