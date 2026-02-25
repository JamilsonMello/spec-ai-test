import { UserDocument } from '../../entity/user';

export interface UserRepository {
  getUserByEmail(email: string): Promise<UserDocument | null>;
  createUser(userData: any): Promise<UserDocument>;
  storePasswordRecoveryToken(userId: string, token: string): Promise<void>;
  getUserByPasswordRecoveryToken(token: string): Promise<UserDocument | null>;
  updatePassword(userId: string, hashedPassword: string): Promise<void>;
  invalidatePasswordRecoveryToken(userId: string): Promise<void>;
  listUsers(filters: { name?: string, email?: string }, pagination: { skip: number, limit: number }, sort: { [key: string]: 'asc' | 'desc' }): Promise<UserDocument[]>;
}
