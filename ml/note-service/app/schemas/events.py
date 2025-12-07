from pydantic import BaseModel, Field
from typing import List
from enum import Enum


class RecordingsPayload(BaseModel):
    class_id: str = Field(..., alias="classId", description="ID klasy/przedmiotu")
    meeting_id: str = Field(..., alias="meetingId", description="ID spotkania/pokoju")
    merged: str = Field(..., description="Ścieżka S3 do pliku merged (np. classId/meetingId.ogg)")
    voices: List[str] = Field(..., description="Lista ścieżek S3 do poszczególnych głosów")

    class Config:
        populate_by_name = True


class DeleteAudioPayload(BaseModel):
    file_paths: List[str] = Field(..., description="Lista ścieżek do usunięcia")


class ProcessingStatus(str, Enum):
    SUCCESS = "SUCCESS"
    FAILURE = "FAILURE"


class ResourcesGeneratedData(BaseModel):
    class_id: str = Field(..., alias="classId", description="ID klasy")
    meeting_id: str = Field(..., alias="meetingId", description="ID spotkania")
    status: ProcessingStatus = Field(..., description="Status przetwarzania")

    class Config:
        populate_by_name = True


class NotificationEvent(BaseModel):
    pattern: str = Field(..., description="Nazwa zdarzenia (np. ResourcesGenerated)")
    data: ResourcesGeneratedData = Field(..., description="Payload zdarzenia")

    class Config:
        populate_by_name = True
