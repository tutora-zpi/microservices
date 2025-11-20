from pydantic import BaseModel, Field
from typing import List


class RecordingsPayload(BaseModel):
    merged: str = Field(
        ...,
        description="Zmergeowane audio pliku"
    )
    voices: List[str] = Field(
        ...,
        description="Pojedyncze głosy"
    )


class DeleteAudioPayload(BaseModel):
    file_paths: List[str] = Field(
        ...,
        description="Lista ścieżek do plików, które należy usunąć"
    )
