export interface UserRepository {
  getUserByEmail(email: string): Promise<any>;
  storePasswordRecoveryToken(userId: string, token: string): Promise<void>;
  getUserIdByPasswordRecoveryToken(token: string): Promise<string | null>;
  updatePassword(userId: string, hashedPassword: string): Promise<void>;
  invalidatePasswordRecoveryToken(userId: string): Promise<void>;
  listUsers(filters: { name?: string, email?: string }, pagination: { skip: number, limit: number }, sort: { [key: string]: 'asc' | 'desc' }): Promise<any[]>;
}