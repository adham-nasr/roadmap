import { Injectable } from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import { Roadmap, RoadmapDocument } from './roadmap.schema';
import { Model } from 'mongoose';

@Injectable()
export class RoadmapRepository {
  constructor(
    @InjectModel(Roadmap.name)
    private readonly roadmapModel: Model<RoadmapDocument>,
  ) {}

  async findAll() {
    return await this.roadmapModel.find({});
  }

  async findOne(id: string) {
    return await this.roadmapModel.findById(id);
  }
}
