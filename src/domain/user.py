import re
import uuid
from datetime import datetime

class InvalidUserError(Exception):
    pass

class User:
    def __init__(self, name: str, email: str, password: str, id: str = None, created_at: datetime = None, updated_at: datetime = None):
        self._id = id if id else str(uuid.uuid4())
        self._created_at = created_at if created_at else datetime.now()
        self._updated_at = updated_at if updated_at else datetime.now()
        self.name = name
        self.email = email
        self.password = password

    @property
    def id(self) -> str:
        return self._id

    @property
    def created_at(self) -> datetime:
        return self._created_at

    @property
    def updated_at(self) -> datetime:
        return self._updated_at

    @property
    def name(self) -> str:
        return self._name

    @name.setter
    def name(self, value: str):
        if not value or not value.strip():
            raise InvalidUserError("Name cannot be empty.")
        self._name = value

    @property
    def email(self) -> str:
        return self._email

    @email.setter
    def email(self, value: str):
        if not re.match(r"[^@]+@[^@]+\.[^@]+", value):
            raise InvalidUserError("Invalid email format.")
        self._email = value

    @property
    def password(self) -> str:
        return self._password

    @password.setter
    def password(self, value: str):
        if len(value) < 8:
            raise InvalidUserError("Password must be at least 8 characters long.")
        self._password = value

    def to_dict(self):
        return {
            "id": self.id,
            "name": self.name,
            "email": self.email,
            "password": self.password,
            "created_at": self.created_at.isoformat(),
            "updated_at": self.updated_at.isoformat(),
        }

    @staticmethod
    def from_dict(data: dict):
        return User(
            id=data.get("id"),
            name=data.get("name"),
            email=data.get("email"),
            password=data.get("password"),
            created_at=datetime.fromisoformat(data["created_at"]) if "created_at" in data else None,
            updated_at=datetime.fromisoformat(data["updated_at"]) if "updated_at" in data else None,
        )
