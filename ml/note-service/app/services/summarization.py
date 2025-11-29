import logging
import os

from app.schemas.downloads import FileType
from app.services.storage_s3 import StorageS3
from app.services.ai_processor import AIProcessor

logger = logging.getLogger(__name__)

NOTES_PROMPT = """
Jesteś profesjonalnym asystentem edukacyjnym. 
Twoim zadaniem jest przygotowanie przejrzystych notatek z lekcji na podstawie podanej transkrypcji.

WYMAGANIA:
1. Użyj formatu Markdown (H1 dla tytułu, H2 dla sekcji, pogrubienia dla kluczowych pojęć).
2. Wypunktuj najważniejsze definicje i wnioski.
3. Zignoruj dygresje i rozmowy niezwiązane z tematem.
4. Styl ma być zwięzły i akademicki.

BARDZO WAŻNE:
Zwróć WYŁĄCZNIE treść notatki w formacie Markdown. 
Nie dodawaj żadnych wstępów typu "Oto notatki", "Jasne, zrobię to" ani zakończeń. 
Zacznij od razu od nagłówka tytułu.
"""

TEST_PROMPT = """
Jesteś asystentem nauczyciela.
Twoim zadaniem jest przygotowanie materiałów sprawdzających wiedzę na podstawie podanej transkrypcji.

WYMAGANIA:
1. Przygotuj 5 pytań zamkniętych (jednokrotnego wyboru A, B, C, D).
2. Przygotuj 2 pytania otwarte wymagające krótszej wypowiedzi.
3. Na samym końcu dokumentu zamieść sekcję "--- KLUCZ ODPOWIEDZI ---" z poprawnymi odpowiedziami.
4. Użyj formatowania Markdown.

BARDZO WAŻNE:
Zwróć WYŁĄCZNIE treść testu i odpowiedzi. 
Nie dodawaj żadnych wstępów typu "Oto test", "Przeanalizowałem tekst" ani innych komentarzy.
Zacznij od razu od treści pytań.
"""

REDUCE_PROMPT = """
Jesteś starszym redaktorem. Poniżej znajdują się fragmenty analizy długiego tekstu.
Scal je w jeden spójny dokument końcowy, zachowując formatowanie Markdown.
Nie dodawaj żadnych komentarzy od siebie. Zwróć tylko scaloną treść.
"""


class SummarizationService:
    def __init__(self, storage_service: StorageS3, ai_processor: AIProcessor):
        self.storage = storage_service
        self.ai_processor = ai_processor

        self.token_limit = 1_000_000
        self.chunk_size_token = 50_000

        logger.info("SummarizationService zainicjalizowany.")

    def generate_and_save_outputs(self, transcript: str, class_id: str, meeting_id: str):
        filename = f"{meeting_id}.md"

        logger.info(f"Rozpoczynam generowanie materiałów dla: {filename}")

        self._generate_notes(transcript, class_id, filename)
        self._generate_tests(transcript, class_id, filename)

        logger.info(f"Zakończono proces generowania dla {filename}.")

    def _generate_notes(self, transcript: str, class_id: str, filename: str):
        logger.info("Generowanie notatek dla ucznia...")

        notes_content = self._process_with_ai(transcript, NOTES_PROMPT)

        if notes_content:
            self.storage.upload_text_result(
                content=notes_content,
                class_id=class_id,
                folder_name=FileType.student_notes.value,
                file_name=filename
            )
        else:
            logger.warning("Wygenerowano puste notatki. Pomijam zapis.")

    def _generate_tests(self, transcript: str, class_id: str, filename: str):
        logger.info("Generowanie testu dla nauczyciela...")

        test_content = self._process_with_ai(transcript, TEST_PROMPT)

        if test_content:
            self.storage.upload_text_result(
                content=test_content,
                class_id=class_id,
                folder_name=FileType.teacher_tests.value,
                file_name=filename
            )
        else:
            logger.warning("Wygenerowano pusty test. Pomijam zapis.")

    def _process_with_ai(self, text: str, system_prompt: str) -> str:
        token_count = self.ai_processor.get_token_count(text)

        if token_count < self.token_limit:
            logger.debug(f"Tekst mieści się w limicie ({token_count} tokenów). Przetwarzanie standardowe.")
            return self.ai_processor.summarize_chunk(text, system_prompt)
        else:
            logger.info(f"Tekst przekracza limit ({token_count} tokenów). Uruchamianie MapReduce.")
            return self._run_map_reduce(text, system_prompt)

    def _run_map_reduce(self, text: str, map_prompt: str) -> str:
        chunks = self._split_text_into_chunks(text)
        partial_results = []

        for i, chunk in enumerate(chunks):
            logger.info(f"MapReduce: Przetwarzanie fragmentu {i + 1}/{len(chunks)}...")
            result = self.ai_processor.summarize_chunk(chunk, map_prompt)
            partial_results.append(result)

        combined_text = "\n\n".join(partial_results)
        logger.info("MapReduce: Scalanie wyników (Reduce step)...")

        return self.ai_processor.summarize_chunk(combined_text, REDUCE_PROMPT)

    def _split_text_into_chunks(self, text: str) -> list[str]:
        char_limit = self.chunk_size_token * 4
        chunks = []
        current_chunk = ""

        for sentence in text.split('. '):
            if len(current_chunk) + len(sentence) < char_limit:
                current_chunk += sentence + ". "
            else:
                chunks.append(current_chunk)
                current_chunk = sentence + ". "

        if current_chunk:
            chunks.append(current_chunk)

        return chunks
