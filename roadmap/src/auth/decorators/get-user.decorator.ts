import { createParamDecorator, ExecutionContext } from '@nestjs/common';
import { Request } from 'express';
import { AuthenticatedUser } from '../interfaces/authenticated-user-interface';

type RequestWithUser = Request & {
  user: AuthenticatedUser;
};

export const GetUser = createParamDecorator(
  (data: keyof AuthenticatedUser | undefined, ctx: ExecutionContext) => {
    const request = ctx.switchToHttp().getRequest<RequestWithUser>();
    const user = request.user;

    // If we pass a property, return only that property
    // Example: @GetUser('id') userId: string
    if (data) {
      return user[data];
    }

    // If we do not pass anything, return the whole user
    // Example: @GetUser() user: AuthenticatedUser
    return user;
  },
);
