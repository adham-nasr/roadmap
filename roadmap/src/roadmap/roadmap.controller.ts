import { Controller, Get, Param } from '@nestjs/common';
import { RoadmapService } from './roadmap.service';
import { ResponseRoadmapDto } from './dto/roadmap_response.dto';
import { ResponseTopicDto } from 'src/topic/dto/topic_response.dto';
import { responseSerializer } from 'src/common/customDecorators/serializer.response';

@Controller('roadmaps')
export class RoadmapController {
  constructor(private readonly roadmapService: RoadmapService) {}

  @Get()
  @responseSerializer(ResponseRoadmapDto)
  async findAll(): Promise<ResponseRoadmapDto[]> {
    const docs = await this.roadmapService.findAll();
    return docs as unknown as ResponseRoadmapDto[];
  }

  @responseSerializer(ResponseRoadmapDto)
  @Get(':id')
  async findOne(@Param('id') id: string) {
    return (await this.roadmapService.findOne(
      id,
    )) as unknown as ResponseRoadmapDto;
  }

  @responseSerializer(ResponseTopicDto)
  @Get(':id/topics')
  async findTopicsByRoadmap(@Param('id') id: string) {
    return (await this.roadmapService.findTopicsByRoadmap(
      id,
    )) as unknown as ResponseTopicDto[];
  }
}
