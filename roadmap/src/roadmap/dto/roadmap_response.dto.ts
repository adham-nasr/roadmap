import { Expose } from 'class-transformer';

export class ResponseRoadmapDto {
  @Expose({ name: '_id' })
  id: string;

  @Expose()
  name: string;

  @Expose()
  description?: string;
}
