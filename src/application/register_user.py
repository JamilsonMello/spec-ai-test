from domain.user import User, InvalidUserError
from domain.user_repository import UserRepository

class RegisterUser:
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository

    def execute(self, name: str, email: str, password: str) -> User:
        if self.user_repository.find_by_email(email):
            raise InvalidUserError("User with this email already exists.")

        user = User(name=name, email=email, password=password)
        self.user_repository.save(user)
        return user
