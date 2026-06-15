import { Expose } from "class-transformer"

export class ResponseRoadmapDto {
    
    @Expose({name:"_id"})
    id_db:string
    @Expose()
    id:string
    @Expose()
    name:string
    @Expose()
    description:string
}
