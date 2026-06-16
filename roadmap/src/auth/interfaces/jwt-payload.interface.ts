import { UserRole } from '../../users/enums/user-role.enum';
// Data stored inside the JWT access token.

export interface JwtPayload {
  sub: string;
  email: string;
  username: string;
  role: UserRole;
}
