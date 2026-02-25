from src.application.use_case import UseCase
from src.application.dtos import RegisterUserInputDTO, RegisterUserOutputDTO
from src.domain.user import User
from src.domain.user_repository import UserRepository

class RegisterUserUseCase(UseCase):
    def __init__(self, user_repository: UserRepository):
        self.user_repository = user_repository

    def execute(self, input_dto: RegisterUserInputDTO) -> RegisterUserOutputDTO:
        if self.user_repository.find_by_email(input_dto.email):
            raise ValueError("User with this email already exists")

        # For simplicity, using a basic hash. In a real application, use bcrypt/argon2.
        # Example: hashed_password = bcrypt.hashpw(input_dto.password.encode('utf-8'), bcrypt.gensalt()).decode('utf-8')
        hashed_password = f"hashed_{input_dto.password}"

        user = User.create_new_user(
            username=input_dto.username,
            email=input_dto.email,
            password_hash=hashed_password
        )

        self.user_repository.save(user)

        return RegisterUserOutputDTO(
            user_id=user.id,
            username=user.username,
            email=user.email
        )
