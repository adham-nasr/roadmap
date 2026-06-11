import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import mongoose, { HydratedDocument } from "mongoose";

@Schema({'collection':'Topic'})
export class Topic{

    @Prop()
    name:string;

    @Prop()
    description:string;

    @Prop()
    type:string;

    @Prop()
    x_axis:number;

    @Prop()
    y_axis:number;

    @Prop()
    repoTopicid:string; /// Sure ?
    
    @Prop({ type: mongoose.Schema.Types.ObjectId, ref: 'Roadmap' })
    roadmap_id: string | mongoose.Types.ObjectId;

    @Prop({ type: mongoose.Schema.Types.ObjectId, ref: 'Topic' })
    parent_topic_id:string | mongoose.Types.ObjectId;

}

export type TopicDocument = HydratedDocument<Topic>

export const topicSchema = SchemaFactory.createForClass(Topic) 