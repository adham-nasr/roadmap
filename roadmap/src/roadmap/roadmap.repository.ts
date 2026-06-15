import { Injectable } from '@nestjs/common';
import { CreateRoadmapDto } from './dto/create-roadmap.dto';
import { UpdateRoadmapDto } from './dto/update-roadmap.dto';
import { InjectModel } from '@nestjs/mongoose';
import { Roadmap, RoadmapDocument } from './roadmap.schema';
import { Model } from 'mongoose';
import { Topic, TopicDocument } from 'src/topic/topic.schema';

@Injectable()
export class RoadmapRepository {
  
  constructor(@InjectModel(Roadmap.name) private readonly roadmapModel:Model<RoadmapDocument>
  ){}
//   create(createRoadmapDto: CreateRoadmapDto) {
//     return 'This action adds a new roadmap';
//   }

  async findAll() {
    return await this.roadmapModel.find({})
  }

  async findOne(id: string) {
    return await this.roadmapModel.findById(id);
  }

//   update(id : string, updateRoadmapDto: UpdateRoadmapDto) {
//     return `This action updates a #${id} roadmap`;
//   }

//   async remove(id : string) {
//     return await this.roadmapModel.deleteOne({_id:id});
//   }

//   async findTopicsByRoadmap(roadmap_id:string){
//     return await this.topicModel.find({roadmap_id:roadmap_id})
//   }
}
