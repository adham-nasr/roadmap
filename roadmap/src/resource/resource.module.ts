import { Module } from '@nestjs/common';
import { ResourceService } from './resource.service';
import { ResourceController } from './resource.controller';
import { MongooseModule } from '@nestjs/mongoose';
import { Resource, resourceSchema } from './resource.schema';

@Module({
  imports:[
      MongooseModule.forFeature([{name:Resource.name,schema:resourceSchema}]) 
    ],
  controllers: [ResourceController],
  providers: [ResourceService],
  exports:[MongooseModule]
})
export class ResourceModule {}
