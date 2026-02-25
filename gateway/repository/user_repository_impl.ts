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

  async getUserIdByPasswordRecoveryToken(token: string): Promise<string | null> {
    try {
      // TODO: Implement getting the user id by password recovery token from the database
            const user = await User.findOne({ passwordRecoveryToken: token });
            if(user){
                console.log(`Getting user id ${user.id} by password recovery token ${token}`);
                return user.id
            }
      console.log(`Getting user id by password recovery token ${token}`);
      return null;
    } catch (error) {
      console.error('Error getting user id by password recovery token:', error);
      return null;
    }
  }

  async updatePassword(userId: string, hashedPassword: string): Promise<void> {
    try {
            const user = await User.findOne({ _id: userId });
            if(user){
                user.password = hashedPassword;
                await user.save();
        console.log(`Updating password for user ${userId} with hashed password ${hashedPassword}`);
            }
    } catch (error) {
      console.error('Error updating password:', error);
    }
  }

  async invalidatePasswordRecoveryToken(userId: string): Promise<void> {
    try {
            const user = await User.findOne({ _id: userId });
            if(user){
                user.passwordRecoveryToken = null;
                await user.save();
        console.log(`Invalidating password recovery token for user ${userId}`);
            }
    } catch (error) {
      console.error('Error invalidating password recovery token:', error);
    }
  }

  async listUsers(filters: { name?: string, email?: string }, pagination: { skip: number, limit: number }, sort: { [key: string]: 'asc' | 'desc' }): Promise<any[]> {
    try {
      const query: any = {};
      if (filters.name) {
        query.name = { $regex: filters.name, $options: 'i' }; // Case-insensitive partial match
      }
      if (filters.email) {
        query.email = { $regex: filters.email, $options: 'i' }; // Case-insensitive partial match
      }

      const users = await User.find(query)
        .sort(sort)
        .skip(pagination.skip)
        .limit(pagination.limit)
        .select('-password -passwordRecoveryToken') // Exclude sensitive fields
        .exec();

      return users;
    } catch (error) {
      console.error('Error listing users:', error);
      return [];
    }
  }
}
