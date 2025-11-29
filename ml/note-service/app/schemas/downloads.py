from enum import Enum
from pydantic import BaseModel
from datetime import datetime
from typing import List


class FileType(str, Enum):
    student_notes = "student_notes"
    teacher_tests = "teacher_tests"
    recording = "recording"


class PresignedUrlResponse(BaseModel):
    url: str
    file_type: FileType
    file_id: str


class FileItem(BaseModel):
    file_id: str
    file_type: FileType
    created_at: datetime
    size_bytes: int


class FileListResponse(BaseModel):
    files: List[FileItem]
