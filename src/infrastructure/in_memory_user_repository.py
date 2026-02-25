from typing import Optional
from domain.user import User
from domain.user_repository import UserRepository

class InMemoryUserRepository(UserRepository):
    def __init__(self):
        self.users = {}

    def save(self, user: User) -> None:
        self.users[user.id] = user

    def find_by_email(self, email: str) -> Optional[User]:
        for user in self.users.values():
            if user.email == email:
                return user
        return None

    def find_by_id(self, user_id: str) -> Optional[User]:
        return self.users.get(user_id)
