from .celery_app import celery_app
from celery.signals import worker_process_init
from app.schemas.events import RecordingsUploaded
from app.services.storage import StorageS3
from app.services.ai_processor import AIProcessor
from app.services.transcription import TranscriptionService
from app.services.summarization import SummarizationService

transcription_service = None
summarization_service = None


@worker_process_init.connect
def on_worker_init(**kwargs):
    """
    Ta funkcja jest automatycznie wywoływana przez Celery raz,
    gdy proces workera jest inicjalizowany.
    To idealne miejsce na załadowanie ciężkich modeli AI.
    """
    print("Sygnał 'worker_process_init': Inicjalizacja serwisów...")

    global transcription_service, summarization_service

    storage_service = StorageS3()
    ai_processor = AIProcessor()

    transcription_service = TranscriptionService(
        storage_service=storage_service,
        ai_processor=ai_processor
    )
    summarization_service = SummarizationService(
        storage_service=storage_service,
        ai_processor=ai_processor
    )

    print("Sygnał 'worker_process_init': Serwisy gotowe do pracy.")


@celery_app.task(name="process_audio_file")
def process_audio_file_task(event_data: dict):
    """
    Główne zadanie Celery, które przyjmuje surowe dane zdarzenia z kolejki,
    waliduje je i deleguje do odpowiednich serwisów.
    """

    print(f"Otrzymano nowe zadanie z danymi: {event_data}")

    if not transcription_service or not summarization_service:
        error_msg = "Serwisy nie zostały zainicjalizowane. Zadanie nie może być wykonane."
        print(error_msg)
        return {"status": "error", "message": error_msg}

    try:
        event = RecordingsUploaded.model_validate(event_data)
    except Exception as e:
        print(f"Błąd walidacji danych: {e}")
        return {"status": "error", "message": "Invalid event data"}

    try:
        transcript = transcription_service.process_recording(
            bucket=event.bucket_name,
            key=event.object_key
        )

        print("\n--- OTRZYMANA TRANSKRYPCJA ---")
        print(transcript)
        print("------------------------------\n")

        summarization_service.generate_and_save_notes(
            transcript=transcript,
            original_filename=event.object_key
        )

        final_message = f"Proces dla pliku {event.object_key} zakończony. Transkrypcja i notatki zostały zapisane."
        print(final_message)
        return {"status": "success", "message": final_message}

    except Exception as e:
        print(f"Zadanie nie powiodło się dla pliku {event.object_key}. Błąd: {e}")
        return {"status": "error", "message": str(e)}
