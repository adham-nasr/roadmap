import {
  ConflictException,
  Injectable,
  InternalServerErrorException,
} from '@nestjs/common';
import { InjectModel } from '@nestjs/mongoose';
import * as bcrypt from 'bcrypt';
import { Model } from 'mongoose';

import { User, UserDocument } from './schemas/user.schema';

type CreateUserData = {
  email: string;
  username: string;
  password: string;
};

@Injectable()
export class UsersService {
  constructor(
    @InjectModel(User.name)
    private readonly userModel: Model<UserDocument>,
  ) {}

  async createUser(createUserData: CreateUserData) {
    const { email, username, password } = createUserData;

    // Hash password before saving it
    const saltRounds = 10;
    const hashedPassword = await bcrypt.hash(password, saltRounds);

    try {
      const createdUser = await this.userModel.create({
        email,
        username,
        password: hashedPassword,
      });

      // Never return password to the client
      return {
        id: createdUser._id.toString(),
        email: createdUser.email,
        username: createdUser.username,
      };
    } catch (error: unknown) {
      if (isDuplicateKeyError(error)) {
        throw new ConflictException('Email or username already exists');
      }

      throw new InternalServerErrorException('Something went wrong');
    }
  }
  async findById(id: string) {
    // Password is not returned because password has select: false in the schema
    return this.userModel.findById(id).exec();
  }

  async findByEmailWithPassword(email: string) {
    // Password is hidden by default, so we must explicitly select it for signin
    return this.userModel.findOne({ email }).select('+password').exec();
  }
}

function isDuplicateKeyError(error: unknown): error is { code: number } {
  return (
    typeof error === 'object' &&
    error !== null &&
    'code' in error &&
    error.code === 11000
  );
}
