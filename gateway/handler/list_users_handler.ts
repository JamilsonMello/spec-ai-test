import { Request, Response } from 'express';
import { ListUsersUseCase } from '../../use_case/list_users_use_case';
import { UserRepositoryImpl } from '../repository/user_repository_impl';

export const listUsersHandler = async (req: Request, res: Response) => {
  try {
    const { name, email, limit, offset } = req.query;

    const filters: { name?: string, email?: string } = {};
    if (typeof name === 'string') {
      filters.name = name;
    }
    if (typeof email === 'string') {
      filters.email = email;
    }

    const pagination = {
      skip: parseInt(offset as string) || 0,
      limit: parseInt(limit as string) || 30, // Default limit to 30 as per spec
    };
    if (pagination.limit > 30) {
        pagination.limit = 30; // Enforce max limit of 30
    }

    const sort = { createdAt: 'desc' }; // Default sort by created_at desc as per spec

    const userRepository = new UserRepositoryImpl();
    const listUsersUseCase = new ListUsersUseCase(userRepository);

    const users = await listUsersUseCase.execute(filters, pagination, sort);

    res.status(200).json(users);
  } catch (error) {
    console.error(error);
    res.status(500).json({ error: 'Internal server error.' });
  }
};
