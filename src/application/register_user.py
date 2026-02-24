from src.application.use_case import UseCase
from src.domain.user import User
from src.domain.user_repository import UserRepository

class RegisterUserRequest:
    def __init__(self, username, email):
        self.username = username
        self.email = email

class RegisterUser(UseCase):
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository

    def execute(self, request: RegisterUserRequest) -> User:
        user = User(id='unique_user_id', username=request.username, email=request.email)
        return self.user_repository.create(user)