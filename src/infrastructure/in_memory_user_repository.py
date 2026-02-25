from typing import Optional, List, Dict

from src.domain.user import User
from src.domain.user_repository import UserRepository

class InMemoryUserRepository(UserRepository):
    def __init__(self, users: Optional[Dict[str, User]] = None):
        self.users = users if users is not None else {}

    def save(self, user: User) -> None:
        self.users[user.id] = user

    def find_by_id(self, user_id: str) -> Optional[User]:
        return self.users.get(user_id)

    def find_by_email(self, email: str) -> Optional[User]:
        for user in self.users.values():
            if user.email == email:
                return user
        return None

    def list_all(self) -> List[User]:
        return list(self.users.values())

    def delete(self, user_id: str) -> None:
        if user_id in self.users:
            del self.users[user_id]
