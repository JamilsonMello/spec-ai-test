// register_use_case.ts

import { UserRepositoryImpl } from '../../gateway/repository/user_repository_impl';
import { UserDTO } from '../../dto/user';
import bcrypt from 'bcrypt';
import jwt from 'jsonwebtoken';
import { UserEntity } from '../../entity/user';
import { ValidationException } from '../../src/exception/validation_exception';

export class RegisterUseCase {
  private userRepository: UserRepositoryImpl;

  constructor(userRepository: UserRepositoryImpl) {
    this.userRepository = userRepository;
  }

  async execute(userDTO: UserDTO): Promise<any> {
    try {
      const newUserEntity = new UserEntity(userDTO.email, userDTO.name, userDTO.password);

      // 1. Validate that the email provided is unique within the database.
      const existingUser = await this.userRepository.getUserByEmail(newUserEntity.email);
      if (existingUser) {
        throw new Error('Email already in use');
      }

      // 2. Validate that the password and password_confirmation fields are identical.
      if (userDTO.password !== userDTO.passwordConfirmation) {
        throw new Error('Passwords do not match');
      }

      // 3. Upon successful validation, the system must persist the user (name, email, hashed password).
      const hashedPassword = await bcrypt.hash(newUserEntity.password, 10);

      const userData = {
        name: newUserEntity.name,
        email: newUserEntity.email,
        password: hashedPassword,
      };

      const newUser = await this.userRepository.createUser(userData);

      // 4. The system must generate and send a welcome email to the user's registered email address, including a unique confirmation link.
      // TODO: Implement email sending logic here

      // 5. The API response for a successful registration must include: `id`, `name`, `email`, `token`, `created_at`, and `updated_at`.
      // 6. The system must ensure that the user's password is never returned in any API response.
      const token = jwt.sign({ userId: newUser._id }, 'secret', { expiresIn: '1h' });

      return {
        id: newUser._id,
        name: newUser.name,
        email: newUser.email,
        token: token,
        createdAt: newUser.createdAt,
        updatedAt: newUser.updatedAt,
      };
    } catch (error) {
      if (error instanceof ValidationException) {
        throw new Error(`Validation Error: ${error.message}`);
      }
      throw error;
    }
  }
}