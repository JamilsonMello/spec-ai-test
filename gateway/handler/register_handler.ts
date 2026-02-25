// register_handler.ts

import { Request, Response } from 'express';
import { RegisterUseCase } from '../../use_case/register_use_case';
import { UserRepositoryImpl } from '../../gateway/repository/user_repository_impl';

export const registerHandler = async (req: Request, res: Response) => {
  try {
    // 1. Extract user data from the request body.
    const { name, email, password, passwordConfirmation } = req.body;

    if (!name || !email || !password || !passwordConfirmation) {
      return res.status(400).json({ error: 'Missing required fields' });
    }

    if (typeof name !== 'string' || name.length < 3) {
      return res.status(400).json({ error: 'Name must be a string with at least 3 characters' });
    }

    if (typeof email !== 'string' || !email.includes('@')) {
      return res.status(400).json({ error: 'Invalid email format' });
    }

    if (typeof password !== 'string' || password.length < 8) {
      return res.status(400).json({ error: 'Password must be a string with at least 8 characters' });
    }

    // 2. Create a UserDTO.
    const userDTO = {
      name,
      email,
      password,
      passwordConfirmation,
    };

    // 3. Instantiate the UserRepository and RegisterUseCase.
    const userRepository = new UserRepositoryImpl();
    const registerUseCase = new RegisterUseCase(userRepository);

    // 4. Execute the RegisterUseCase.
    const newUser = await registerUseCase.execute(userDTO);

    // 5. Send the response.
    res.status(201).json(newUser);
  } catch (error) {
    console.error(error);
    res.status(400).json({ error: error.message });
  }
};