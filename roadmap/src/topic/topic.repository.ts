import { Injectable } from '@nestjs/common';

import { InjectModel } from '@nestjs/mongoose';
import { Topic, TopicDocument } from './topic.schema';
import { Model } from 'mongoose';

@Injectable()
export class TopicRepository {
  constructor(
    @InjectModel(Topic.name) private readonly topicModel: Model<TopicDocument>,
  ) {}

  async findAll() {
    return await this.topicModel.find();
  }

  async findOne(id: string) {
    return await this.topicModel.find({ topicId: id });
  }

  async findByRoadmap(roadmap_id: string) {
    return await this.topicModel.find({ roadmapid: roadmap_id });
  }
}
