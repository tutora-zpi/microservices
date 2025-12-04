import os
from app.services.storage_s3 import StorageS3
from app.services.ai_processor import AIProcessor


class TranscriptionService:
    def __init__(self, storage_service: StorageS3, ai_processor: AIProcessor):
        self.storage = storage_service
        self.ai_processor = ai_processor
        print("TranscriptionService zainicjalizowany.")

    def process_recording(self, bucket: str, key: str) -> str:
        local_audio_path = None
        try:
            print(f"Rozpoczynanie pobierania pliku: s3://{bucket}/{key}")
            local_audio_path = self.storage.download_audio(object_name=key)

            if not local_audio_path:
                raise FileNotFoundError(f"Nie udało się pobrać pliku {key} z bucketa {bucket}.")

            print(f"Plik pobrany do: {local_audio_path}. Przekazywanie do transkrypcji...")
            transcript = self.ai_processor.transcribe(audio_path=local_audio_path)

            print(f"Transkrypcja dla {key} zakończona pomyślnie.")

            return transcript

        except Exception as e:
            print(f"Wystąpił krytyczny błąd podczas procesu transkrypcji dla {key}: {e}")
            raise

        finally:
            if local_audio_path and os.path.exists(local_audio_path):
                os.remove(local_audio_path)
                print(f"Posprzątano plik tymczasowy: {local_audio_path}")