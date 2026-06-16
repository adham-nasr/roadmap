import { Injectable } from '@nestjs/common';
import { AuthGuard } from '@nestjs/passport';

@Injectable()
// Protects routes by validating the JWT access token using JwtStrategy.
export class JwtAuthGuard extends AuthGuard('jwt') {}
