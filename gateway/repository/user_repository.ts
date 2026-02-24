export interface UserRepository {
  getUserByEmail(email: string): Promise<any>;
}