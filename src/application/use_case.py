from abc import ABC, abstractmethod

class UseCase(ABC):
    @abstractmethod
    def execute(self, input_dto: any) -> any:
        pass
