import { Injectable } from '@nestjs/common';
import { TopicRepository } from './topic.repository';

@Injectable()
export class TopicService {


  constructor(private readonly topicRepo: TopicRepository
){}
  // create(createTopicDto: CreateTopicDto) {
  //   return 'This action adds a new topic';
  // }

  async findAll() {
    return await this.topicRepo.findAll()
  }

  // async findAllByRoadmap(roadmap_id:string) {
  //   return await this.topicModel.find({roadmap_id:roadmap_id})
  // }
  async findOne(id : string) {
    return await this.topicRepo.findOne(id);
  }

  // update(id : string, updateTopicDto: UpdateTopicDto) {
  //   return `This action updates a #${id} topic`;
  // }

  // remove(id : string) {
  //   return `This action removes a #${id} topic`;
  // }

  // async findResourcesByTopic(topic_id:string){
  //   return await this.resourceModel.find({topic_id:topic_id})
  // }
}
