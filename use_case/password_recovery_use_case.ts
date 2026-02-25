import { UserRepository } from '../domain/user_repository';
import { v4 as uuidv4 } from 'uuid';

export class PasswordRecoveryUseCase {
  constructor(private userRepository: UserRepository) {}

  async execute(email: string): Promise<void> {
    // Check if the email exists in the database
    const user = await this.userRepository.getUserByEmail(email);

    // If the email doesn't exist, return a success message to avoid exposing user data
    if (!user) {
      return;
    }

    // Generate a unique token
    const token = uuidv4();

    // Store the token in the database
    await this.userRepository.storePasswordRecoveryToken(user.id, token);

    // Send the email with the recovery link
    this.sendPasswordRecoveryEmail(email, token);
  }

  private sendPasswordRecoveryEmail(email: string, token: string): void {
    // TODO: Implement email sending logic
    const recoveryLink = `http://localhost:3000/reset-password?token=${token}`;
    console.log(`Sending password recovery email to ${email} with recovery link: ${recoveryLink}`);
  }
}