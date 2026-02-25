import express, { Request, Response } from 'express';
import { PasswordRecoveryUseCase } from '../../use_case/password_recovery_use_case';
import { UserRepositoryImpl } from '../../gateway/repository/user_repository_impl';

const userRepository = new UserRepositoryImpl();
const passwordRecoveryUseCase = new PasswordRecoveryUseCase(userRepository);

export const passwordRecoveryHandler = async (req: Request, res: Response) => {
  const { email } = req.body;


  try {
    await passwordRecoveryUseCase.execute(email);
    return res.status(200).send({ message: 'Instructions sent to your email' });
  } catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' });
  }
};