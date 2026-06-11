import { Controller, Get, Post, Body, Patch, Param, Delete } from '@nestjs/common';
import { RoadmapService } from './roadmap.service';
import { CreateRoadmapDto } from './dto/create-roadmap.dto';
import { UpdateRoadmapDto } from './dto/update-roadmap.dto';
import { TopicService } from 'src/topic/topic.service';

@Controller('roadmaps')
export class RoadmapController {
  constructor(private readonly roadmapService: RoadmapService) {}

  @Post()
  create(@Body() createRoadmapDto: CreateRoadmapDto) {
    return this.roadmapService.create(createRoadmapDto);
  }

  @Get()
  findAll() {
    return this.roadmapService.findAll();
  }

  @Get(':id')
  findOne(@Param('id') id: string) {
    return this.roadmapService.findOne(id);
  }

  @Patch(':id')
  update(@Param('id') id: string, @Body() updateRoadmapDto: UpdateRoadmapDto) {
    return this.roadmapService.update(id, updateRoadmapDto);
  }

  @Delete(':id')
  remove(@Param('id') id: string) {
    return this.roadmapService.remove(id);
  }


  @Get(':id/topics')
  async findTopicsByRoadmap(@Param('id') id:string){
    return await this.roadmapService.findTopicsByRoadmap(id);
  }
}
