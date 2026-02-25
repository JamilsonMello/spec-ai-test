from http import HTTPStatus
from src.application.dtos import (
    RegisterUserInputDTO,
    PasswordRecoveryInputDTO,
    ListUsersOutputDTO,
    UserOutputDTO
)
from src.application.register_user import RegisterUserUseCase
from src.application.password_recovery_use_case import PasswordRecoveryUseCase
from src.application.list_users_use_case import ListUsersUseCase


class RegisterUserController:
    def __init__(self, register_user_use_case: RegisterUserUseCase):
        self.register_user_use_case = register_user_use_case

    def handle(self, request_data: dict):
        try:
            input_dto = RegisterUserInputDTO(
                username=request_data.get("username"),
                email=request_data.get("email"),
                password=request_data.get("password")
            )
            output_dto = self.register_user_use_case.execute(input_dto)
            return {
                "status": HTTPStatus.CREATED,
                "body": {
                    "id": output_dto.user_id,
                    "username": output_dto.username,
                    "email": output_dto.email
                }
            }
        except ValueError as e:
            return {"status": HTTPStatus.BAD_REQUEST, "body": {"message": str(e)}}
        except Exception as e:
            return {"status": HTTPStatus.INTERNAL_SERVER_ERROR, "body": {"message": "An unexpected error occurred"}}


class PasswordRecoveryController:
    def __init__(self, password_recovery_use_case: PasswordRecoveryUseCase):
        self.password_recovery_use_case = password_recovery_use_case

    def handle(self, request_data: dict):
        try:
            input_dto = PasswordRecoveryInputDTO(
                email=request_data.get("email")
            )
            output_dto = self.password_recovery_use_case.execute(input_dto)
            return {"status": HTTPStatus.OK, "body": {"message": output_dto.message}}
        except ValueError as e:
            return {"status": HTTPStatus.BAD_REQUEST, "body": {"message": str(e)}}
        except Exception as e:
            return {"status": HTTPStatus.INTERNAL_SERVER_ERROR, "body": {"message": "An unexpected error occurred"}}


class ListUsersController:
    def __init__(self, list_users_use_case: ListUsersUseCase):
        self.list_users_use_case = list_users_use_case

    def handle(self):
        try:
            output_dto = self.list_users_use_case.execute(None)
            users_data = [
                {"id": user.id, "username": user.username, "email": user.email}
                for user in output_dto.users
            ]
            return {"status": HTTPStatus.OK, "body": {"users": users_data}}
        except Exception as e:
            return {"status": HTTPStatus.INTERNAL_SERVER_ERROR, "body": {"message": "An unexpected error occurred"}}
