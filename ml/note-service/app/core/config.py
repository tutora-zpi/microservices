from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    # --- RABBIT ---
    RABBITMQ_BROKER_URL: str

    # --- AWS / S3 ---
    AWS_ACCESS_KEY_ID: str
    AWS_SECRET_ACCESS_KEY: str
    AWS_REGION: str
    S3_RECORDINGS_BUCKET_NAME: str
    S3_NOTES_BUCKET_NAME: str

    # --- AI: Whisper (Transkrypcja) ---
    WHISPER_MODEL_ID: str = "openai/whisper-base"
    WHISPER_BATCH_SIZE: int = 8
    WHISPER_CHUNK_LENGTH: int = 30

    # --- AI: Gemini (LLM) ---
    GEMINI_API_KEY: str
    GEMINI_MODEL_NAME: str = "gemini-1.5-pro"

    # --- System / Storage ---
    LOCAL_TMP_DIR: str = "/tmp/transcriptions"

    class Config:
        env_file = ".env"


settings = Settings()
