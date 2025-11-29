import logging
from app.services.ai.transcriber import LocalTranscriber
from app.services.ai.gemini_client import GeminiClient

logger = logging.getLogger(__name__)


class AIProcessor:
    def __init__(self):
        logger.info("AIProcessor: Inicjalizacja podsystemów...")

        self.transcriber = LocalTranscriber()
        self.llm_client = GeminiClient()

        logger.info("AIProcessor: Wszystkie podsystemy gotowe.")

    def transcribe(self, audio_path: str) -> str:
        return self.transcriber.transcribe(audio_path)

    def summarize_chunk(self, chunk: str, prompt: str) -> str:
        full_prompt = f"{prompt}\n\nTranskrypcja:\n---\n{chunk}\n---\n"
        return self.llm_client.generate_content(full_prompt)

    def get_token_count(self, text: str) -> int:
        return self.llm_client.count_tokens(text)