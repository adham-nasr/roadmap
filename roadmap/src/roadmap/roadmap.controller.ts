import { Controller, Get, Post, Body, Patch, Param, Delete, UseInterceptors } from '@nestjs/common';
import { RoadmapService } from './roadmap.service';
import { TopicService } from 'src/topic/topic.service';
import { ResponseRoadmapDto } from './dto/roadmap_response.dto';
import { ResponseTopicDto } from 'src/topic/dto/topic_response.dto';
import { responseSerializer } from 'src/common/customDecorators/serializer.response';

@Controller('roadmaps')
export class RoadmapController {
  constructor(private readonly roadmapService: RoadmapService) {}

  // @Post()
  // create(@Body() createRoadmapDto: CreateRoadmapDto) {
  //   return this.roadmapService.create(createRoadmapDto);
  // }

  @responseSerializer(ResponseRoadmapDto)
  @Get()
  findAll() {
    return this.roadmapService.findAll();
  }
  
  @responseSerializer(ResponseRoadmapDto)
  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.roadmapService.findOne(id);
  }

  // @Patch(':id')
  // update(@Param('id') id: string, @Body() updateRoadmapDto: UpdateRoadmapDto) {
  //   return this.roadmapService.update(id, updateRoadmapDto);
  // }

  // @Delete(':id')
  // remove(@Param('id') id: string) {
  //   return this.roadmapService.remove(id);
  // }


  @responseSerializer(ResponseTopicDto)
  @Get(':id/topics')
  async findTopicsByRoadmap(@Param('id') id:string){
    return await this.roadmapService.findTopicsByRoadmap(id);
  }
}
