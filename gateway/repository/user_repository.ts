export interface UserRepository {
  getUserByEmail(email: string): Promise<any>;
  storePasswordRecoveryToken(userId: string, token: string): Promise<void>;
  getUserIdByPasswordRecoveryToken(token: string): Promise<string | null>;
  updatePassword(userId: string, hashedPassword: string): Promise<void>;
  invalidatePasswordRecoveryToken(userId: string): Promise<void>;
}