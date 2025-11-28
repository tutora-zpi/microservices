# app/services/storage_s3.py
import os
import logging
import boto3
from botocore.exceptions import ClientError
from typing import Optional, TYPE_CHECKING, List, Dict, Any

if TYPE_CHECKING:
    from app.core.config import Settings

from app.core.config import settings as global_settings

logger = logging.getLogger(__name__)


class StorageS3:
    def __init__(self, settings: Optional['Settings'] = None):
        self.config = settings or global_settings

        self.region = self.config.AWS_REGION
        self.recordings_bucket = self.config.S3_RECORDINGS_BUCKET_NAME
        self.notes_bucket = self.config.S3_NOTES_BUCKET_NAME
        self.local_tmp_dir = self.config.LOCAL_TMP_DIR

        self.s3_client = boto3.client(
            's3',
            region_name=self.region,
            aws_access_key_id=self.config.AWS_ACCESS_KEY_ID,
            aws_secret_access_key=self.config.AWS_SECRET_ACCESS_KEY,
            aws_session_token=self.config.AWS_SESSION_TOKEN
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

    def generate_presigned_url(self, object_key: str, bucket_type: str = "notes", expiration: int = 3600) -> str:
        if bucket_type == "notes":
            bucket_name = self.notes_bucket
        elif bucket_type == "recordings":
            bucket_name = self.recordings_bucket
        else:
            logger.error(f"Nieznany typ bucketa: {bucket_type}")
            return ""

        try:
            url = self.s3_client.generate_presigned_url(
                'get_object',
                Params={'Bucket': bucket_name, 'Key': object_key},
                ExpiresIn=expiration
            )
            return url
        except ClientError as e:
            logger.error(f"Błąd generowania URL dla {object_key}: {e}")
            return ""

    def upload_text_result(self, content: str, class_id: str, folder_name: str, file_name: str):
        s3_key = f"{class_id}/{folder_name}/{file_name}"

        safe_folder_name = folder_name.replace("/", "_")
        local_temp_filename = f"{safe_folder_name}_{file_name}"
        local_path = os.path.join(self.local_tmp_dir, local_temp_filename)

        try:
            with open(local_path, 'w', encoding='utf-8') as f:
                f.write(content)

            logger.info(f"Wysyłanie: {local_path} -> s3://{self.notes_bucket}/{s3_key}")

            self.s3_client.upload_file(
                Filename=local_path,
                Bucket=self.notes_bucket,
                Key=s3_key,
                ExtraArgs={'ContentType': 'text/markdown; charset=utf-8'}
            )

            if os.path.exists(local_path):
                os.remove(local_path)

        except Exception as e:
            logger.error(f"Błąd podczas uploadu wyniku do {s3_key}: {e}", exc_info=True)
            if os.path.exists(local_path):
                os.remove(local_path)
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

    def list_files(self, prefix: str) -> List[Dict[str, Any]]:
        """
        Listuje pliki w buckecie NOTES zaczynające się od danego prefiksu (np. 'student_notes/').
        Zwraca listę słowników z kluczami, datami i rozmiarem.
        """
        try:
            if not prefix.endswith('/'):
                prefix += '/'

            response = self.s3_client.list_objects_v2(
                Bucket=self.notes_bucket,
                Prefix=prefix
            )

            files = []
            if 'Contents' in response:
                for obj in response['Contents']:
                    if obj['Key'] == prefix:
                        continue

                    files.append({
                        "key": obj['Key'],
                        "last_modified": obj['LastModified'],
                        "size": obj['Size']
                    })
            return files

        except ClientError as e:
            logger.error(f"Błąd listowania plików w {prefix}: {e}")
            return []
