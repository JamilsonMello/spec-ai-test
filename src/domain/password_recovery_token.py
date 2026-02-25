import uuid
from datetime import datetime, timedelta

class InvalidPasswordRecoveryTokenError(Exception):
    pass

class PasswordRecoveryToken:
    def __init__(self, token: str, user_id: str, expires_at: datetime, id: str = None, created_at: datetime = None):
        self._id = id if id else str(uuid.uuid4())
        self._created_at = created_at if created_at else datetime.now()
        self.token = token
        self.user_id = user_id
        self.expires_at = expires_at

    @property
    def id(self) -> str:
        return self._id

    @property
    def created_at(self) -> datetime:
        return self._created_at

    @property
    def token(self) -> str:
        return self._token

    @token.setter
    def token(self, value: str):
        if not value or not value.strip():
            raise InvalidPasswordRecoveryTokenError("Token cannot be empty.")
        self._token = value

    @property
    def user_id(self) -> str:
        return self._user_id

    @user_id.setter
    def user_id(self, value: str):
        if not value or not value.strip():
            raise InvalidPasswordRecoveryTokenError("User ID cannot be empty.")
        self._user_id = value

    @property
    def expires_at(self) -> datetime:
        return self._expires_at

    @expires_at.setter
    def expires_at(self, value: datetime):
        if value <= datetime.now():
            raise InvalidPasswordRecoveryTokenError("Expiration date must be in the future.")
        self._expires_at = value

    def is_expired(self) -> bool:
        return datetime.now() >= self.expires_at

    @classmethod
    def create(cls, user_id: str):
        token = str(uuid.uuid4())
        expires_at = datetime.now() + timedelta(hours=24)
        return cls(token=token, user_id=user_id, expires_at=expires_at)

    def to_dict(self):
        return {
            "id": self.id,
            "token": self.token,
            "user_id": self.user_id,
            "expires_at": self.expires_at.isoformat(),
            "created_at": self.created_at.isoformat(),
        }

    @staticmethod
    def from_dict(data: dict):
        return PasswordRecoveryToken(
            id=data.get("id"),
            token=data.get("token"),
            user_id=data.get("user_id"),
            expires_at=datetime.fromisoformat(data["expires_at"]),
            created_at=datetime.fromisoformat(data["created_at"]) if "created_at" in data else None,
        )
