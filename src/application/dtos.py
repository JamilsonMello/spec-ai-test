from dataclasses import dataclass
from typing import List, Optional

@dataclass
class RegisterUserInputDTO:
    username: str
    email: str
    password: str

@dataclass
class RegisterUserOutputDTO:
    user_id: str
    username: str
    email: str

@dataclass
class PasswordRecoveryInputDTO:
    email: str

@dataclass
class PasswordRecoveryOutputDTO:
    message: str # e.g., "Password recovery email sent"

@dataclass
class UserOutputDTO:
    id: str
    username: str
    email: str

@dataclass
class ListUsersOutputDTO:
    users: List[UserOutputDTO]

@dataclass
class PasswordResetInputDTO:
    token: str
    new_password: str

@dataclass
class PasswordResetOutputDTO:
    message: str
