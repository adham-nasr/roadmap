import { Controller, Get, Param, Query } from '@nestjs/common';
import { TopicService } from './topic.service';
import { responseSerializer } from 'src/common/customDecorators/serializer.response';
import { ResponseTopicDto } from './dto/topic_response.dto';

@Controller('topics')
@responseSerializer(ResponseTopicDto)
export class TopicController {
  constructor(private readonly topicService: TopicService) {}

  @Get()
  async findAll() {
    const res = await this.topicService.findAll();
    return res as unknown as ResponseTopicDto[];
  }

  @Get(':id')
  async findOne(@Param('id') id: string) {
    return (await this.topicService.findOne(id)) as unknown as ResponseTopicDto;
  }

  @Get()
  async findByRoadmapId(@Query('roadmapId') roadmapId: string) {
    const res = await this.topicService.findByRoadmapId(roadmapId);
    return res as unknown as ResponseTopicDto[];
  }
}
