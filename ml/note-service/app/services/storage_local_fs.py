import os
import shutil
from pathlib import Path


class StorageLocalFS:
    def __init__(self, base_path: str = "local_storage"):
        self.base_path = Path(base_path).resolve()
        self.input_path = self.base_path / "input"
        self.output_path = self.base_path / "output"
        self.transcript_path = self.base_path / "transcript"

        self.input_path.mkdir(parents=True, exist_ok=True)
        self.output_path.mkdir(parents=True, exist_ok=True)
        self.transcript_path.mkdir(parents=True, exist_ok=True)

        print(f"StorageLocalFS (mockup) zainicjalizowany. Input: '{self.input_path}', Output: '{self.output_path}'")

    def download_audio(self, object_name: str) -> str | None:
        local_file_path = self.input_path / object_name
        print(f"Symulacja pobierania: sprawdzanie ścieżki '{local_file_path}'...")

        if not local_file_path.exists():
            print(f"Błąd: Plik {local_file_path} nie istnieje w lokalnym magazynie.")
            return None

        print(f"Plik '{object_name}' znaleziony lokalnie.")
        return str(local_file_path)

    def upload_notes(self, local_file_path: str, object_name: str) -> bool:
        destination_path = self.output_path / object_name
        print(f"Symulacja wysyłania: kopiowanie '{local_file_path}' do '{destination_path}'...")
        try:
            shutil.copy(local_file_path, destination_path)
            print("Kopiowanie zakończone pomyślnie.")
            return True
        except Exception as e:
            print(f"Błąd podczas kopiowania pliku: {e}")
            return False

    def delete_audio(self, object_name: str):
        """
        Symuluje 'usuwanie' pliku poprzez usunięcie go z katalogu 'local_storage'.
        """
        file_to_delete = self.input_path / object_name
        print(f"Symulacja usuwania: usuwanie pliku '{file_to_delete}'...")
        try:
            if file_to_delete.exists():
                os.remove(file_to_delete)
                print("Usuwanie pliku zakończone pomyślnie.")
            else:
                print("Plik już nie istnieje, pomijanie usuwania.")
        except Exception as e:
            print(f"Błąd podczas usuwania pliku: {e}")

    def save_debug_file(self, folder: str, filename: str, content: str):
        if folder == 'transcript':
            destination_path = self.transcript_path / filename
        else:
            print(f"Nieznany folder debugowania: {folder}")
            return

        print(f"Zapisywanie pliku debugowania do: {destination_path}")
        try:
            with open(destination_path, 'w', encoding='utf-8') as f:
                f.write(content)
        except Exception as e:
            print(f"Błąd podczas zapisywania pliku debugowania: {e}")