import torch
from transformers import pipeline
import google.generativeai as genai
from app.core.config import settings


class AIProcessor:
    def __init__(self):
        device = "cuda:0" if torch.cuda.is_available() else "cpu"
        print(f"AIProcessor: Inicjalizacja na urządzeniu: {device}")
        try:
            self.transcriber = pipeline(
                "automatic-speech-recognition",
                model="openai/whisper-base",
                device=device
            )
            print("AIProcessor: Model Whisper został pomyślnie załadowany.")

        except Exception as e:
            print(f"AIProcessor: Krytyczny błąd podczas ładowania modelu: {e}")
            raise

        try:
            print("AIProcessor: Konfiguracja klienta API Gemini...")
            genai.configure(api_key=settings.GEMINI_API_KEY)

            safety_settings = [
                {"category": "HARM_CATEGORY_HARASSMENT", "threshold": "BLOCK_NONE"},
                {"category": "HARM_CATEGORY_HATE_SPEECH", "threshold": "BLOCK_NONE"},
                {"category": "HARM_CATEGORY_SEXUALLY_EXPLICIT", "threshold": "BLOCK_NONE"},
                {"category": "HARM_CATEGORY_DANGEROUS_CONTENT", "threshold": "BLOCK_NONE"},
            ]

            self.gemini_model = genai.GenerativeModel(
                'gemini-2.5-pro',
                safety_settings=safety_settings
            )
            print("AIProcessor: Klient API Gemini gotowy.")
        except Exception as e:
            print(f"AIProcessor: Krytyczny błąd podczas konfiguracji API Gemini: {e}")
            raise

    def transcribe(self, audio_path: str) -> str:
        if not audio_path:
            print("AIProcessor: Otrzymano pustą ścieżkę do pliku audio.")
            raise ValueError("Ścieżka do pliku audio nie może być pusta.")

        print(f"AIProcessor: Rozpoczynanie transkrypcji pliku: {audio_path}...")

        try:
            outputs = self.transcriber(
                audio_path,
                return_timestamps=True
            )

            transcript = outputs['text'].strip()
            print("AIProcessor: Transkrypcja zakończona.")

            return transcript
        except Exception as e:
            print(f"AIProcessor: Wystąpił błąd podczas transkrypcji pliku {audio_path}: {e}")
            raise

    def summarize_chunk(self, chunk: str, prompt: str) -> str:
        print("AIProcessor: Wysyłanie fragmentu do API Gemini...")
        try:
            full_prompt = f"{prompt}\n\nTranskrypcja:\n---\n{chunk}\n---\n"
            response = self.gemini_model.generate_content(full_prompt)
            print("AIProcessor: Otrzymano odpowiedź z API Gemini.")
            return response.text
        except Exception as e:
            print(f"AIProcessor: Błąd podczas wywołania API Gemini: {e}")
            return ""

    def get_token_count(self, text: str) -> int:
        return self.gemini_model.count_tokens(text).total_tokens
