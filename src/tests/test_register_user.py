import unittest
from unittest.mock import Mock
from src.application.register_user import RegisterUser, RegisterUserRequest
from src.domain.user import User
from src.domain.user_repository import UserRepository

class TestRegisterUser(unittest.TestCase):
    def test_register_user(self):
        # Mock the UserRepository
        user_repository = Mock(UserRepository)
        user_repository.create.return_value = User(id='test_id', username='test_user', email='test@example.com')

        # Create the use case with the mock repository
        register_user = RegisterUser(user_repository)

        # Create a request
        request = RegisterUserRequest(username='test_user', email='test@example.com')

        # Execute the use case
        user = register_user.execute(request)

        # Assertions
        self.assertEqual(user.username, 'test_user')
        self.assertEqual(user.email, 'test@example.com')
        user_repository.create.assert_called_once()

if __name__ == '__main__':
    unittest.main()