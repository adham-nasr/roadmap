import { Injectable } from '@nestjs/common';
import { CreateResourceDto } from './dto/create-resource.dto';
import { UpdateResourceDto } from './dto/update-resource.dto';
import { InjectModel } from '@nestjs/mongoose';
import { Resource, ResourceDocument } from './resource.schema';
import { Model } from 'mongoose';

@Injectable()
export class ResourceService {

  constructor(@InjectModel(Resource.name) private readonly resourceModel:Model<ResourceDocument>){}
  create(createResourceDto: CreateResourceDto) {
    return 'This action adds a new resource';
  }

  async findAll() {
    return await this.resourceModel.find({});
  }

  // findAllForTopic(topic_id:string) {
  //   return this.resourceModel.find({});
  // }

  async findOne(id : string) {
    return await this.resourceModel.findById(id);
  }

  update(id : string, updateResourceDto: UpdateResourceDto) {
    return `This action updates a #${id} resource`;
  }

  remove(id : string) {
    return `This action removes a #${id} resource`;
  }
}
