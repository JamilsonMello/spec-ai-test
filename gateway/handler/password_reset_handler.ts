import express, { Request, Response } from 'express';
import { PasswordResetUseCase } from '../../use_case/password_reset_use_case';
import { UserRepositoryImpl } from '../../gateway/repository/user_repository_impl';

const userRepository = new UserRepositoryImpl();
const passwordResetUseCase = new PasswordResetUseCase(userRepository);

export const passwordResetHandler = async (req: Request, res: Response) => {
  const { token, password, confirmPassword } = req.body;

  // FR-005: The password and confirmation password must be identical
  if (password !== confirmPassword) {
    return res.status(400).send({ message: 'Passwords do not match' });
  }

  try {
    await passwordResetUseCase.execute(token, password);
    return res.status(200).send({ message: 'Password reset successfully' });
  } catch (error) {
    console.error(error);
    return res.status(500).send({ message: 'Internal server error' });
  }
};