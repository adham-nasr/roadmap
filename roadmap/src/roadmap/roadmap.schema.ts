import { Prop, Schema, SchemaFactory } from '@nestjs/mongoose';
import { HydratedDocument } from 'mongoose';

@Schema({ collection: 'roadmaps' })
export class Roadmap {
  @Prop({ required: true })
  name: string;
  @Prop()
  description?: string;
}

export type RoadmapDocument = HydratedDocument<Roadmap>;

export const roadmapSchema = SchemaFactory.createForClass(Roadmap);
