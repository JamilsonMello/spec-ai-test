// user_repository_impl.ts

import mongoose from 'mongoose';
import { UserSchema } from '../../entities/user';
import { UserRepository } from './user_repository';

const User = mongoose.model('User', UserSchema);

export class UserRepositoryImpl implements UserRepository {
  async getUserByEmail(email: string): Promise<any | null> {
    try {
      const user = await User.findOne({ email });
      return user;
    } catch (error) {
      console.error('Error getting user by email:', error);
      return null;
    }
  }

  async createUser(userData: any): Promise<any | null> {
    try {
      const user = new User(userData);
      await user.save();
      return user;
    } catch (error) {
      console.error('Error creating user:', error);
      return null;
    }
  }
}