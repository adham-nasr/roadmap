import { Expose, Type } from 'class-transformer';

export class ChildTopicDto {
  @Expose()
  targetId: string;
  @Expose()
  relation: string;
}

export class ResourcesDto {
  @Expose()
  type: string;

  @Expose()
  title: string;

  @Expose()
  link: string;
}

export class ResponseTopicDto {
  @Expose()
  topicId: string;
  @Expose()
  name: string;
  @Expose()
  description: string;
  @Expose()
  type: string;
  @Expose()
  position: {
    x: number;
    y: number;
  };
  @Expose()
  roadmapid: string;

  @Expose()
  @Type(() => ChildTopicDto)
  childTopics: ChildTopicDto[];

  @Type(() => ResourcesDto)
  @Expose()
  resources: ResourcesDto[];
}
