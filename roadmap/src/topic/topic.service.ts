import { Injectable } from '@nestjs/common';
import { CreateTopicDto } from './dto/create-topic.dto';
import { UpdateTopicDto } from './dto/update-topic.dto';
import { InjectModel } from '@nestjs/mongoose';
import { Topic, TopicDocument } from './topic.schema';
import { Model } from 'mongoose';
import { Resource, ResourceDocument } from 'src/resource/resource.schema';

@Injectable()
export class TopicService {

  constructor(@InjectModel(Topic.name) private readonly topicModel:Model<TopicDocument> ,
              @InjectModel(Resource.name) private readonly resourceModel:Model<ResourceDocument>
){}
  create(createTopicDto: CreateTopicDto) {
    return 'This action adds a new topic';
  }

  async findAll() {
    return await this.topicModel.find()
  }

  // async findAllByRoadmap(roadmap_id:string) {
  //   return await this.topicModel.find({roadmap_id:roadmap_id})
  // }
  async findOne(id : string) {
    return await this.topicModel.findById(id);
  }

  update(id : string, updateTopicDto: UpdateTopicDto) {
    return `This action updates a #${id} topic`;
  }

  remove(id : string) {
    return `This action removes a #${id} topic`;
  }

  async findResourcesByTopic(topic_id:string){
    return await this.resourceModel.find({topic_id:topic_id})
  }
}
