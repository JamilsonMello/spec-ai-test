from typing import List

from src.application.use_case import UseCase
from src.application.dtos import ListUsersOutputDTO, UserOutputDTO
from src.domain.user_repository import UserRepository

class ListUsersUseCase(UseCase):
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository

    def execute(self, input_dto: None) -> ListUsersOutputDTO:
        users = self.user_repository.list_all()
        user_dtos = [UserOutputDTO(id=user.id, username=user.username, email=user.email) for user in users]
        return ListUsersOutputDTO(users=user_dtos)
