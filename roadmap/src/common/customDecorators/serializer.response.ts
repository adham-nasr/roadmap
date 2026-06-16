import { UseInterceptors } from '@nestjs/common';
import { SerializerInterceptor } from '../interceptors/serializer';

export function responseSerializer(dto: any) {
  return UseInterceptors(new SerializerInterceptor(dto));
}
