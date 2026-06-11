import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import mongoose, { HydratedDocument } from "mongoose";



@Schema({collection:'Resource'})
export class Resource{
    
    @Prop()
    link:string;

    @Prop()
    type:string;

    @Prop({type:mongoose.Schema.Types.ObjectId , ref:'Topic'})
    topic_id: string | mongoose.Types.ObjectId;

}

export type ResourceDocument = HydratedDocument<Resource>;

export const resourceSchema = SchemaFactory.createForClass(Resource);

