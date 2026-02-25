from src.domain.user_repository import UserRepository
from src.infrastructure.in_memory_user_repository import InMemoryUserRepository

from src.application.register_user import RegisterUserUseCase
from src.application.password_recovery_use_case import PasswordRecoveryUseCase
from src.application.list_users_use_case import ListUsersUseCase
from src.application.dtos import RegisterUserInputDTO, PasswordRecoveryInputDTO

from src.presentation.controllers import (
    RegisterUserController,
    PasswordRecoveryController,
    ListUsersController
)

def setup_dependencies():
    """Sets up and returns all dependencies."""
    user_repository: UserRepository = InMemoryUserRepository()

    register_user_use_case = RegisterUserUseCase(user_repository)
    password_recovery_use_case = PasswordRecoveryUseCase(user_repository)
    list_users_use_case = ListUsersUseCase(user_repository)

    register_user_controller = RegisterUserController(register_user_use_case)
    password_recovery_controller = PasswordRecoveryController(password_recovery_use_case)
    list_users_controller = ListUsersController(list_users_use_case)

    return {
        "user_repository": user_repository, # For initial data loading demonstration
        "register_user_controller": register_user_controller,
        "password_recovery_controller": password_recovery_controller,
        "list_users_controller": list_users_controller,
    }

def main():
    dependencies = setup_dependencies()
    register_user_controller = dependencies["register_user_controller"]
    password_recovery_controller = dependencies["password_recovery_controller"]
    list_users_controller = dependencies["list_users_controller"]
    user_repository = dependencies["user_repository"] # For pre-populating users for listing

    print("--- Simulating User Registration ---")
    register_data_1 = {"username": "testuser1", "email": "test1@example.com", "password": "password123"}
    response_1 = register_user_controller.handle(register_data_1)
    print(f"Register User 1 Response: {response_1}")

    register_data_2 = {"username": "testuser2", "email": "test2@example.com", "password": "password456"}
    response_2 = register_user_controller.handle(register_data_2)
    print(f"Register User 2 Response: {response_2}")

    print("\n--- Simulating Duplicate User Registration ---")
    response_duplicate = register_user_controller.handle(register_data_1)
    print(f"Register Duplicate User Response: {response_duplicate}")

    print("\n--- Simulating User Listing ---")
    list_response = list_users_controller.handle()
    print(f"List Users Response: {list_response}")

    print("\n--- Simulating Password Recovery ---")
    recovery_data = {"email": "test1@example.com"}
    recovery_response = password_recovery_controller.handle(recovery_data)
    print(f"Password Recovery Response: {recovery_response}")

    print("\n--- Simulating Password Recovery for Non-Existent User ---")
    non_existent_recovery_data = {"email": "nonexistent@example.com"}
    non_existent_recovery_response = password_recovery_controller.handle(non_existent_recovery_data)
    print(f"Non-Existent Password Recovery Response: {non_existent_recovery_response}")


if __name__ == "__main__":
    main()
