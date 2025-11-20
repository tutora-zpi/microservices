import logging
import torch
from transformers import pipeline
from app.core.config import settings

logger = logging.getLogger(__name__)


class LocalTranscriber:
    def __init__(self):
        self.device = "cuda:0" if torch.cuda.is_available() else "cpu"
        self.model_id = settings.WHISPER_MODEL_ID
        self.pipeline = None

        self._load_model()

    def _load_model(self):
        logger.info(f"Inicjalizacja modelu Whisper ({self.model_id}) na urządzeniu: {self.device}")
        try:
            torch_dtype = torch.float16 if self.device != "cpu" else torch.float32

            self.pipeline = pipeline(
                "automatic-speech-recognition",
                model=self.model_id,
                device=self.device,
                torch_dtype=torch_dtype
            )
            logger.info("Model Whisper został pomyślnie załadowany.")
        except Exception as e:
            logger.critical(f"Krytyczny błąd ładowania modelu Whisper: {e}", exc_info=True)
            raise

    def transcribe(self, audio_path: str) -> str:
        if not audio_path:
            raise ValueError("Ścieżka do pliku audio jest pusta.")

        logger.info(f"Start transkrypcji pliku: {audio_path}")
        try:
            outputs = self.pipeline(
                audio_path,
                chunk_length_s=settings.WHISPER_CHUNK_LENGTH,
                batch_size=settings.WHISPER_BATCH_SIZE,
                return_timestamps=True
            )

            text = outputs['text'].strip()
            logger.info(f"Transkrypcja zakończona. Wygenerowano {len(text)} znaków.")
            return text

        except Exception as e:
            logger.error(f"Błąd podczas transkrypcji {audio_path}: {e}", exc_info=True)
            raise
