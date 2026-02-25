from abc import ABC, abstractmethod
from typing import Optional, List
from src.domain.user import User

class UserRepository(ABC):
    @abstractmethod
    def save(self, user: User) -> None:
        pass

    @abstractmethod
    def find_by_id(self, user_id: str) -> Optional[User]:
        pass

    @abstractmethod
    def find_by_email(self, email: str) -> Optional[User]:
        pass

    @abstractmethod
    def list_all(self) -> List[User]:
        pass

    @abstractmethod
    def delete(self, user_id: str) -> None:
        pass
