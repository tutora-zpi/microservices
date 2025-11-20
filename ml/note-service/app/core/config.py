from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", env_file_encoding='utf-8')

    RABBITMQ_BROKER_URL: str

    AWS_ACCESS_KEY_ID: str
    AWS_SECRET_ACCESS_KEY: str
    AWS_REGION: str
    S3_RECORDINGS_BUCKET_NAME: str
    S3_NOTES_BUCKET_NAME: str

    GEMINI_API_KEY: str


settings = Settings()
