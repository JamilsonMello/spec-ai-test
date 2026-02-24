// register_handler.ts

import { Request, Response } from 'express';
import { RegisterUseCase } from '../../use_case/register_use_case';
import { UserRepositoryImpl } from '../../gateway/repository/user_repository_impl';

export const registerHandler = async (req: Request, res: Response) => {
  try {
    // 1. Extract user data from the request body.
    const { name, email, password, passwordConfirmation } = req.body;

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