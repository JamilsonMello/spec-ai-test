from abc import ABC, abstractmethod
from src.domain.user import User

class UserRepository(ABC):
    @abstractmethod
    def get_by_id(self, user_id: str) -> User:
        pass

    @abstractmethod
    def create(self, user: User) -> User:
        pass