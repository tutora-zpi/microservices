import logging
import google.generativeai as genai
from app.core.config import settings

logger = logging.getLogger(__name__)

SAFETY_SETTINGS = [
    {"category": "HARM_CATEGORY_HARASSMENT", "threshold": "BLOCK_NONE"},
    {"category": "HARM_CATEGORY_HATE_SPEECH", "threshold": "BLOCK_NONE"},
    {"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT", "threshold": "BLOCK_NONE"},
    {"category": "HARM_CATEGORY_DANGEROUS_CONTENT", "threshold": "BLOCK_NONE"},
]


class GeminiClient:
    def __init__(self):
        self.model_name = settings.GEMINI_MODEL_NAME
        self.model = None
        self._configure_api()

    def _configure_api(self):
        logger.info("Konfiguracja klienta Gemini API...")
        try:
            if not settings.GEMINI_API_KEY:
                raise ValueError("Brak GEMINI_API_KEY w ustawieniach.")

            genai.configure(api_key=settings.GEMINI_API_KEY)

            self.model = genai.GenerativeModel(
                self.model_name,
                safety_settings=SAFETY_SETTINGS
            )
            logger.info(f"Klient Gemini gotowy (Model: {self.model_name}).")
        except Exception as e:
            logger.critical(f"Błąd konfiguracji Gemini: {e}", exc_info=True)
            raise

    def generate_content(self, prompt: str) -> str:
        logger.debug("Wysyłanie zapytania do Gemini...")
        try:
            response = self.model.generate_content(prompt)
            if not response.text:
                logger.warning("Gemini zwróciło pustą odpowiedź.")
                return ""

            return response.text
        except Exception as e:
            logger.error(f"Błąd API Gemini: {e}", exc_info=True)
            return ""

    def count_tokens(self, text: str) -> int:
        try:
            return self.model.count_tokens(text).total_tokens
        except Exception as e:
            logger.error(f"Błąd liczenia tokenów: {e}")
            return 0