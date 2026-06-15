import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Roadmap, RoadmapDocument } from './roadmap.schema';
import { Model } from 'mongoose';
import { Topic, TopicDocument } from 'src/topic/topic.schema';
import { RoadmapRepository } from './roadmap.repository';
import { TopicRepository } from 'src/topic/topic.repository';

@Injectable()
export class RoadmapService {
  
  constructor(private readonly roadmapRepo:RoadmapRepository,
    private readonly topicRepo:TopicRepository
  ){}
  // create(createRoadmapDto: CreateRoadmapDto) {
  //   return 'This action adds a new roadmap';
  // }

  async findAll() {
    return await this.roadmapRepo.findAll()
  }

  async findOne(id: string) {
    return await this.roadmapRepo.findOne(id);
  }

  // update(id : string, updateRoadmapDto: UpdateRoadmapDto) {
  //   return `This action updates a #${id} roadmap`;
  // }

  // async remove(id : string) {
  //   return await this.roadmapModel.deleteOne({_id:id});
  // }

  async findTopicsByRoadmap(roadmap_id:string){
    return await this.topicRepo.findByRoadmap(roadmap_id);
  }
}
