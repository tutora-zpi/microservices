# app/services/storage_s3.py
import os
import logging
import boto3
from botocore.exceptions import ClientError
from app.core.config import settings

logger = logging.getLogger(__name__)


class StorageS3:
    def __init__(self):
        self.recordings_bucket = settings.S3_RECORDINGS_BUCKET_NAME
        self.notes_bucket = settings.S3_NOTES_BUCKET_NAME
        self.region = settings.AWS_REGION

        self.local_tmp_dir = settings.LOCAL_TMP_DIR

        self.s3_client = boto3.client(
            's3',
            region_name=self.region,
            aws_access_key_id=settings.AWS_ACCESS_KEY_ID,
            aws_secret_access_key=settings.AWS_SECRET_ACCESS_KEY
        )

    def download_audio(self, object_name: str) -> str:
        try:
            local_file_path = os.path.join(self.local_tmp_dir, object_name)

            os.makedirs(os.path.dirname(local_file_path), exist_ok=True)

            logger.info(f"Rozpoczynanie pobierania z S3: s3://{self.recordings_bucket}/{object_name} -> {local_file_path}")

            self.s3_client.download_file(
                Bucket=self.recordings_bucket,
                Key=object_name,
                Filename=local_file_path
            )

            logger.debug(f"Pobieranie zakończone sukcesem: {local_file_path}")
            return local_file_path

        except ClientError as e:
            logger.error(f"Błąd krytyczny S3 podczas pobierania pliku {object_name}: {e}", exc_info=True)
            raise e

    def upload_notes(self, local_file_path: str, object_key: str):
        try:
            logger.info(f"Wysyłanie notatki: {local_file_path} -> s3://{self.notes_bucket}/{object_key}")

            self.s3_client.upload_file(
                Filename=local_file_path,
                Bucket=self.notes_bucket,
                Key=object_key,
                ExtraArgs={'ContentType': 'text/plain'}
            )
            logger.info(f"Pomyślnie wysłano notatkę do {self.notes_bucket}")

        except ClientError as e:
            logger.error(f"Błąd wysyłania notatki (Bucket: {self.notes_bucket}): {e}", exc_info=True)
            raise e

    def delete_audio(self, object_key: str):
        logger.info(f"Usuwanie pliku audio: s3://{self.recordings_bucket}/{object_key}")
        try:
            self.s3_client.delete_object(
                Bucket=self.recordings_bucket,
                Key=object_key
            )
            logger.info(f"Usunięto plik audio: {object_key}")
        except ClientError as e:
            logger.warning(f"Nie udało się usunąć pliku audio {object_key}: {e}")
