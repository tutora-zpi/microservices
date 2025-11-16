# app/main.py

from fastapi import FastAPI, BackgroundTasks
from app.schemas.events import AudioUploadedEvent
from app.tasks.meeting_processor import process_audio_file_task

app = FastAPI(
    title="Meeting Notes Service",
    description="Mikroserwis do transkrypcji i podsumowywania spotkań."
)


@app.get("/", summary="Endpoint statusu")
def read_root():
    """Sprawdza, czy serwis działa."""
    return {"status": "ok"}


@app.post("/process-audio", summary="Zleć przetwarzanie nagrania")
def trigger_audio_processing(event: AudioUploadedEvent):
    """
    Testowy endpoint do symulowania otrzymania zdarzenia 'AudioUploadedEvent'.
    Przyjmuje dane o pliku i wysyła zadanie do przetworzenia w tle przez Celery.
    """
    print(f"Otrzymano zlecenie przez API dla pliku: {event.object_key}")

    # Wysyłamy zadanie do kolejki RabbitMQ.
    # .delay() to skrót do wysłania zadania.
    # .model_dump() konwertuje obiekt Pydantic na słownik, który Celery może zserializować.
    process_audio_file_task.delay(event.model_dump())

    return {
        "message": "Zadanie przetwarzania audio zostało pomyślnie zlecone.",
        "event_data": event
    }
