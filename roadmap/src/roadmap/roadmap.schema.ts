import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import { HydratedDocument } from "mongoose";



@Schema({'collection':'Roadmap'})
export class Roadmap{

    @Prop()
    name:string;
    @Prop()
    description:string;
}

export type RoadmapDocument = HydratedDocument<Roadmap>

export const roadmapSchema = SchemaFactory.createForClass(Roadmap)