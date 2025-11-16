# app/services/storage_s3.py

import boto3
from botocore.exceptions import ClientError
from app.core.config import settings

class StorageS3:
    def __init__(self):
        # ... (bez zmian)
        self.s3_client = boto3.client(...)
        self.bucket_name = settings.S3_BUCKET_NAME

    def download_audio(self, object_name: str) -> str:
        # ... (bez zmian)
        pass

    def upload_notes(self, local_file_path: str, object_name: str):
        # ... (bez zmian)
        pass

    # === NOWA METODA ===
    def delete_audio(self, object_name: str):
        """Usuwa obiekt (plik) z bucketa S3."""
        print(f"Rozpoczynanie usuwania pliku s3://{self.bucket_name}/{object_name}...")
        try:
            self.s3_client.delete_object(Bucket=self.bucket_name, Key=object_name)
            print("Usuwanie zakończone pomyślnie.")
        except ClientError as e:
            # Logujemy błąd, ale nie przerywamy działania,
            # ponieważ główny cel (transkrypcja) został osiągnięty.
            print(f"Błąd podczas usuwania pliku {object_name} z S3: {e}")