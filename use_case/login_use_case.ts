import { UserRepository } from '../../gateway/repository/user_repository';
import { UserDTO } from '../../dto/user';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';

export class LoginUseCase {
  private userRepository: UserRepository;

  constructor(userRepository: UserRepository) {
    this.userRepository = userRepository;
  }

  async execute(userDTO: UserDTO): Promise<string | null> {
    const user = await this.userRepository.getUserByEmail(userDTO.email);

    if (!user) {
      return null;
    }

    const passwordMatch = await bcrypt.compare(userDTO.password, user.password);

    if (!passwordMatch) {
      return null;
    }

    const token = jwt.sign({ userId: user._id }, 'secret', { expiresIn: '1h' });

    return token;
  }
}