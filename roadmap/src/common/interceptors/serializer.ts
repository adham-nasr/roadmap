import { CallHandler, ExecutionContext, NestInterceptor } from "@nestjs/common";
import { plainToClass, plainToInstance } from "class-transformer";
import { map, Observable } from "rxjs";


export class SerializerInterceptor implements NestInterceptor{

    constructor(private dto:any){}
    intercept(context: ExecutionContext, next: CallHandler<any>): Observable<any> | Promise<Observable<any>> {
        
        //

        return next.handle().pipe(
            map((data:any)=>{
                return plainToInstance(this.dto,data,{
                    excludeExtraneousValues:true
                    // enableCircularCheck:true,
                    // strategy:'exposeAll',
                    // excludePrefixes: ['$'], 
                })
            }),
        )
    }
}