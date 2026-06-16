import { UserRole } from '../../users/enums/user-role.enum';
// User data attached to request.user after successful JWT authentication.

export interface AuthenticatedUser {
  id: string;
  email: string;
  username: string;
  role: UserRole;
}
