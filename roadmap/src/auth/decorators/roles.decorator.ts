import { SetMetadata } from '@nestjs/common';

import { UserRole } from '../../users/enums/user-role.enum';

export const ROLES_KEY = 'roles';

// This decorator lets us write: @Roles(UserRole.ADMIN)
export const Roles = (...roles: UserRole[]) => SetMetadata(ROLES_KEY, roles);
