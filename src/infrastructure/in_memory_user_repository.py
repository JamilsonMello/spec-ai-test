from src.domain.user import User
from src.domain.user_repository import UserRepository

class InMemoryUserRepository(UserRepository):
    def __init__(self):
        self.users = {}

    def get_by_id(self, user_id: str) -> User:
        return self.users.get(user_id)

    def create(self, user: User) -> User:
        self.users[user.id] = user
        return user