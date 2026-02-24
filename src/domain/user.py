from src.domain.entity import Entity

class User(Entity):
    def __init__(self, id, username, email):
        super().__init__(id)
        self.username = username
        self.email = email

    def __str__(self):
        return f"User(id={self.id}, username={self.username}, email={self.email})"