import { Prop, Schema, SchemaFactory } from "@nestjs/mongoose";
import mongoose, { HydratedDocument } from "mongoose";

@Schema({_id:false})
export class Resources{
    @Prop({type: String,required:true})
    type:string

    @Prop({required:true})
    title:string

    @Prop({required:true})
    link:string
}

export const ResourcesSchema = SchemaFactory.createForClass(Resources);


@Schema({_id:false})
export class ChildTopic{
    @Prop({required:true})
    targetId: string

    @Prop({required:true})
    relation:string
 
}

export const ChildTopicSchema = SchemaFactory.createForClass(ChildTopic);


@Schema({'collection':'Topic'})
export class Topic{

    @Prop({required:true})
    name:string;

    @Prop()
    description:string;

    @Prop({type:String , required:true})
    type:string

    @Prop({type:Object})
    position:{
        x:number,
        y:number
    }

    @Prop()
    repoTopicid:string; /// Sure ?
    

    @Prop({ type: mongoose.Schema.Types.ObjectId, ref: 'Roadmap' , required:true , index:true})
    roadmap_id: string | mongoose.Types.ObjectId;
    

    @Prop({ type:[ChildTopicSchema] })
    childTopics: ChildTopic[]

    @Prop({type:[ResourcesSchema]})
    resources: Resources[]

}

export type TopicDocument = HydratedDocument<Topic>

export const topicSchema = SchemaFactory.createForClass(Topic) 