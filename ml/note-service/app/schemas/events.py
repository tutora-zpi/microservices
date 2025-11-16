from pydantic import BaseModel, Field


class AudioUploadedEvent(BaseModel):
    bucket_name: str = Field(
        ...,
        description="Nazwa bucketa S3, w którym znajduje się nagranie.",
        examples=["tutora-meeting-recordings"]
    )
    object_key: str = Field(
        ...,
        description="Klucz (ścieżka do pliku) nagrania w buckecie S3.",
        examples=["raw/meeting-12345.wav"]
    )
