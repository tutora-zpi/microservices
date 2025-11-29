from typing import Union, List
from pydantic_settings import BaseSettings
from pydantic import AnyHttpUrl, field_validator


class Settings(BaseSettings):
    # --- RABBIT ---
    RABBITMQ_BROKER_URL: str

    # --- AWS / S3 ---
    AWS_ACCESS_KEY_ID: str
    AWS_SECRET_ACCESS_KEY: str
    AWS_SESSION_TOKEN: str | None = None
    AWS_REGION: str
    S3_RECORDINGS_BUCKET_NAME: str
    S3_NOTES_BUCKET_NAME: str

    # --- AI: Whisper (Transkrypcja) ---
    WHISPER_MODEL_ID: str = "openai/whisper-base"
    WHISPER_BATCH_SIZE: int = 8
    WHISPER_CHUNK_LENGTH: int = 30

    # --- AI: Gemini (LLM) ---
    GEMINI_API_KEY: str
    GEMINI_MODEL_NAME: str = "gemini-2.5-pro"

    # --- System / Storage ---
    LOCAL_TMP_DIR: str = "/tmp/transcriptions"

    # --- API Security & Config ---
    BACKEND_CORS_ORIGINS: Union[List[str], str] = []

    @field_validator("BACKEND_CORS_ORIGINS", mode="before")
    def assemble_cors_origins(cls, v: Union[str, List[str]]) -> List[str]:
        """
        Parsuje ciąg znaków oddzielony przecinkami na listę stringów.
        Pozwala to w .env wpisać: BACKEND_CORS_ORIGINS="http://localhost:3000,http://localhost:8080"
        """
        if isinstance(v, str) and not v.startswith("["):
            return [i.strip() for i in v.split(",")]
        elif isinstance(v, (list, str)):
            return v
        raise ValueError(v)

    AUTH_ALGORITHM: str = "RS256"
    AUTH_JWKS_URL: str

    class Config:
        env_file = ".env"
        case_sensitive = True


settings = Settings()
