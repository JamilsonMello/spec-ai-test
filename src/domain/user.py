import uuid
import re

class User:
    def __init__(self, id: str, username: str, email: str, password_hash: str):
        if not id:
            raise ValueError("User ID cannot be empty")
        if not username or len(username) < 3:
            raise ValueError("Username must be at least 3 characters long")
        if not self._is_valid_email(email):
            raise ValueError("Invalid email format")
        if not password_hash:
            raise ValueError("Password hash cannot be empty")

        self.id = id
        self.username = username
        self.email = email
        self.password_hash = password_hash

    @staticmethod
    def create_new_user(username: str, email: str, password_hash: str) -> 'User':
        return User(str(uuid.uuid4()), username, email, password_hash)

    def update_username(self, new_username: str):
        if not new_username or len(new_username) < 3:
            raise ValueError("Username must be at least 3 characters long")
        self.username = new_username

    def update_email(self, new_email: str):
        if not self._is_valid_email(new_email):
            raise ValueError("Invalid email format")
        self.email = new_email

    def update_password_hash(self, new_password_hash: str):
        if not new_password_hash:
            raise ValueError("Password hash cannot be empty")
        self.password_hash = new_password_hash

    @staticmethod
    def _is_valid_email(email: str) -> bool:
        # Basic email validation
        return re.match(r"[^@]+@[^@]+\.[^@]+", email) is not None
