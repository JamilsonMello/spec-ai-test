import { UserRepository } from '../domain/user_repository';

export class PasswordResetUseCase {
  constructor(private userRepository: UserRepository) {}

  async execute(token: string, password: string): Promise<void> {
    // Verify the token
    const userId = await this.userRepository.getUserIdByPasswordRecoveryToken(token);

    if (!userId) {
      throw new Error('Invalid or expired token');
    }

    // Hash the new password
    const hashedPassword = await this.hashPassword(password);

    // Update the password in the database
    await this.userRepository.updatePassword(userId, hashedPassword);

    // Invalidate the token
    await this.userRepository.invalidatePasswordRecoveryToken(userId);
  }

  private async hashPassword(password: string): Promise<string> {
    // TODO: Implement password hashing logic
    console.log(`Hashing password ${password}`);
    return password;
  }
}