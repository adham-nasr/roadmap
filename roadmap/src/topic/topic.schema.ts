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
    @Prop({ type: mongoose.Schema.Types.ObjectId, ref: 'Topic' })
    targetId: string | mongoose.Types.ObjectId

    @Prop({required:true})
    relation:string
 
}

export const ChildTopicSchema = SchemaFactory.createForClass(ChildTopic);


@Schema({'collection':'Topic'})
export class Topic{

    @Prop()
    name:string;

    @Prop()
    description:string;

    @Prop({type:String})
    type:string

    @Prop({type:Object})
    position:{
        x:number,
        y:number
    }

    @Prop()
    repoTopicid:string; /// Sure ?
    
    @Prop({ type: mongoose.Schema.Types.ObjectId, ref: 'Roadmap' })
    roadmap_id: string | mongoose.Types.ObjectId;
    

    @Prop({ type:[ChildTopicSchema] })
    ChildTopics: ChildTopic[]

    @Prop({type:[ResourcesSchema]})
    Resources: Resources[]

}

export type TopicDocument = HydratedDocument<Topic>

export const topicSchema = SchemaFactory.createForClass(Topic) 