import os
import tempfile
from .storage_local_fs import StorageLocalFS
from .ai_processor import AIProcessor

MAP_PROMPT = """Jesteś asystentem AI. Twoim zadaniem jest stworzenie zwięzłych notatek z poniższego fragmentu transkrypcji spotkania. 
Skup się wyłącznie na kluczowych decyzjach, przydzielonych zadaniach (action items) i najważniejszych wnioskach.
Ignoruj zwykłe rozmowy i wypełniacze. Użyj formatu Markdown.
"""

REDUCE_PROMPT = """Jesteś starszym redaktorem AI. Poniżej znajduje się seria notatek cząstkowych z długiego spotkania, 
oddzielonych liniami "---". Twoim zadaniem jest zredagowanie ich w jeden, spójny, logicznie 
uporządkowany dokument końcowy w formacie Markdown. 
Usuń redundancje i stwórz profesjonalne podsumowanie całego spotkania.
"""


class SummarizationService:
    def __init__(self, storage_service: StorageLocalFS, ai_processor: AIProcessor):
        self.storage = storage_service
        self.ai_processor = ai_processor
        self.token_limit = 1_000_000
        self.chunk_size_token = 50_000
        print("SummarizationService zainicjalizowany (tryb API Gemini).")

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
        chunks.append(current_chunk)

        print(f"Podzielono tekst na {len(chunks)} fragmentów.")
        return chunks

    def generate_and_save_notes(self, transcript: str, original_filename: str):
        if self.ai_processor.get_token_count(transcript) < self.token_limit:
            print("Transkrypcja mieści się w limicie. Wykonywanie pojedynczego podsumowania.")
            final_notes_content = self.ai_processor.summarize_chunk(transcript, MAP_PROMPT)
        else:
            print("Transkrypcja przekracza limit. Uruchamianie procedury MapReduce...")
            transcript_chunks = self._split_text_into_chunks(transcript)

            partial_summaries = []
            for i, chunk in enumerate(transcript_chunks):
                print(f"Przetwarzanie fragmentu {i + 1}/{len(transcript_chunks)}...")
                partial_summary = self.ai_processor.summarize_chunk(chunk, MAP_PROMPT)
                partial_summaries.append(partial_summary)

            print("Łączenie notatek cząstkowych...")
            combined_summary = "\n\n---\n[Kolejny fragment]\n---\n\n".join(partial_summaries)

            print("Wysyłanie połączonych notatek do finalnej redakcji (Reduce)...")
            final_notes_content = self.ai_processor.summarize_chunk(combined_summary, REDUCE_PROMPT)

        print("\n--- OSTATECZNE NOTATKI (z API Gemini) ---")
        print(final_notes_content)
        print("------------------------------------------\n")

        base_name = os.path.splitext(original_filename)[0]
        notes_filename = f"{base_name}_notes.txt"

        temp_dir = tempfile.gettempdir()
        local_notes_path = os.path.join(temp_dir, notes_filename)
        with open(local_notes_path, 'w', encoding='utf-8') as f:
            f.write(final_notes_content)

        self.storage.upload_notes(
            local_file_path=local_notes_path,
            object_name=notes_filename
        )
        os.remove(local_notes_path)