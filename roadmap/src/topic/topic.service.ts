import { Injectable } from '@nestjs/common';
import { TopicRepository } from './topic.repository';

@Injectable()
export class TopicService {
  constructor(private readonly topicRepo: TopicRepository) {}

  async findAll() {
    return await this.topicRepo.findAll();
  }

  async findOne(id: string) {
    return await this.topicRepo.findOne(id);
  }

  async findByRoadmapId(roadmapId: string) {
    return await this.topicRepo.findByRoadmap(roadmapId);
  }
}
