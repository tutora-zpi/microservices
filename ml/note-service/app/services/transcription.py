import os
from app.services.storage_s3 import StorageS3
from app.services.ai_processor import AIProcessor


class TranscriptionService:
    """
    Serwis orkiestrujący proces transkrypcji nagrania.
    """

    def __init__(self, storage_service: StorageS3, ai_processor: AIProcessor):
        """
        Inicjalizuje serwis, wstrzykując zależności do S3 i procesora AI.
        """
        self.storage = storage_service
        self.ai_processor = ai_processor
        print("TranscriptionService zainicjalizowany.")

    def process_recording(self, bucket: str, key: str) -> str:
        """
        Główna metoda biznesowa: pobiera nagranie, dokonuje transkrypcji i sprząta.

        Args:
            bucket (str): Nazwa bucketa S3.
            key (str): Klucz obiektu (pliku) w S3.

        Returns:
            str: Wynikowa transkrypcja w formie tekstu.
        """
        local_audio_path = None
        try:
            # Krok 1: Pobranie pliku audio z S3 na dysk lokalny workera
            print(f"Rozpoczynanie pobierania pliku: s3://{bucket}/{key}")
            local_audio_path = self.storage.download_audio(object_name=key)

            if not local_audio_path:
                raise FileNotFoundError(f"Nie udało się pobrać pliku {key} z bucketa {bucket}.")

            # Krok 2: Przekazanie pobranego pliku do modelu AI w celu transkrypcji
            print(f"Plik pobrany do: {local_audio_path}. Przekazywanie do transkrypcji...")
            transcript = self.ai_processor.transcribe(audio_path=local_audio_path)

            # base_name = os.path.splitext(key)[0]
            # transcript_filename = f"{base_name}_transcript.txt"
            # self.storage.save_debug_file('transcript', transcript_filename, transcript)

            print(f"Transkrypcja dla {key} zakończona pomyślnie.")
            # self.storage.delete_audio(object_name=key)

            return transcript

        except Exception as e:
            print(f"Wystąpił krytyczny błąd podczas procesu transkrypcji dla {key}: {e}")
            raise

        finally:
            if local_audio_path and os.path.exists(local_audio_path):
                os.remove(local_audio_path)
                print(f"Posprzątano plik tymczasowy: {local_audio_path}")