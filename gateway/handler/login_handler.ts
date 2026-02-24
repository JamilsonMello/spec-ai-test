import express, { Request, Response } from 'express';
import { LoginUseCase } from '../../use_case/login_use_case';
import { InMemoryUserRepository } from '../../src/infrastructure/in_memory_user_repository';
import { UserDTO } from '../../dto/user';

const userRepository = new InMemoryUserRepository();
const loginUseCase = new LoginUseCase(userRepository);

export const loginHandler = async (req: Request, res: Response) => {
  const { email, password } = req.body;

  // FR-003: Validate email format
  if (!email || !/^\w-\.]+@([\w-]+\.)+[\w-]{2,4}$/.test(email)) {
    return res.status(400).send({ message: 'Invalid email format' });
  }

  const userDTO: UserDTO = { email, password };

  try {
    const token = await loginUseCase.execute(userDTO);

    if (!token) {
      return res
        .status(401)
        .send({ message: 'email ou senha inválidos' });
    }

    res.status(200).send({ token });
  } catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' });
  }
};