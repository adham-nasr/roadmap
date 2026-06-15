import { Controller, Get, Post, Body, Patch, Param, Delete } from '@nestjs/common';
import { TopicService } from './topic.service';
import { responseSerializer } from 'src/common/customDecorators/serializer.response';
import { ResponseTopicDto } from './dto/topic_response.dto';

@Controller('topics')
@responseSerializer(ResponseTopicDto)
export class TopicController {
  constructor(private readonly topicService: TopicService) {}

  // @Post()
  // create(@Body() createTopicDto: CreateTopicDto) {
  //   return this.topicService.create(createTopicDto);
  // }

  @Get()
  async findAll(){
    const res = await this.topicService.findAll();
    return res as unknown as ResponseTopicDto[]
  }

  @Get(':id')
  async findOne(@Param('id') id: string) {
    return await this.topicService.findOne(id) as unknown as ResponseTopicDto
  }

  // @Patch(':id')
  // update(@Param('id') id: string, @Body() updateTopicDto: UpdateTopicDto) {
  //   return this.topicService.update(id, updateTopicDto);
  // }

  // @Delete(':id')
  // remove(@Param('id') id: string) {
  //   return this.topicService.remove(id);
  // }

  // @Get(':id/resources')
  // async findResourcesByTopic(@Param('id') id:string) {
  //   return await this.topicService.findResourcesByTopic(id);
  // }

}
