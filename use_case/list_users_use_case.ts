import { UserRepository } from '../gateway/repository/user_repository';

export class ListUsersUseCase {
  private userRepository: UserRepository;

  constructor(userRepository: UserRepository) {
    this.userRepository = userRepository;
  }

  async execute(filters: { name?: string, email?: string }, pagination: { skip: number, limit: number }, sort: { [key: string]: 'asc' | 'desc' }): Promise<any[]> {
    // Business logic for listing users, if any, goes here.
    // For now, it directly calls the repository method.
    return this.userRepository.listUsers(filters, pagination, sort);
  }
}
