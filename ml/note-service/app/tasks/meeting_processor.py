from celery.utils.log import get_task_logger
from celery.signals import worker_process_init
from .celery_app import celery_app

from app.schemas.events import RecordingsPayload, ProcessingStatus
from app.services.storage_s3 import StorageS3
from app.services.ai_processor import AIProcessor
from app.services.transcription import TranscriptionService
from app.services.summarization import SummarizationService
from app.services.notification_publisher import NotificationPublisher
from app.tasks.file_cleaner import delete_audio_files_task
from app.core.config import settings

logger = get_task_logger(__name__)

_transcription_service: TranscriptionService = None
_summarization_service: SummarizationService = None
_notification_publisher: NotificationPublisher = None


@worker_process_init.connect
def on_worker_init(**kwargs):
    global _transcription_service, _summarization_service, _notification_publisher

    logger.info("Inicjalizacja serwisów AI i Storage...")

    try:
        storage = StorageS3()
        ai_processor = AIProcessor()
        _notification_publisher = NotificationPublisher()

        _transcription_service = TranscriptionService(
            storage_service=storage,
            ai_processor=ai_processor
        )
        _summarization_service = SummarizationService(
            storage_service=storage,
            ai_processor=ai_processor
        )
        logger.info("Serwisy zainicjalizowane pomyślnie.")

    except Exception as e:
        logger.critical(f"Błąd krytyczny inicjalizacji serwisów: {e}", exc_info=True)


@celery_app.task(name="process_audio_file")
def process_audio_file_task(event_data: dict):
    if not _check_services_availability():
        return {"status": "error", "message": "Services not initialized"}

    try:
        event = _parse_event_data(event_data)

        transcript = _perform_transcription(event.merged)

        _perform_summarization(transcript, event.class_id, event.meeting_id)

        _trigger_cleanup(event)

        _send_notification(event.class_id, event.meeting_id, ProcessingStatus.SUCCESS)

        logger.info(f"Proces zakończony sukcesem dla: {event.merged}")
        return {"status": "success", "file": event.merged}

    except ValueError as e:
        logger.error(f"Błąd danych wejściowych: {e}")
        return {"status": "error", "message": str(e)}

    except Exception as e:
        logger.error(f"Błąd przetwarzania taska: {e}", exc_info=True)
        if event:
            _send_notification(event.class_id, event.meeting_id, ProcessingStatus.SUCCESS)
        raise e


def _check_services_availability() -> bool:
    if not _transcription_service or not _summarization_service:
        logger.error("Próba uruchomienia zadania bez zainicjalizowanych serwisów.")
        return False
    return True


def _parse_event_data(raw_data: dict) -> RecordingsPayload:
    try:
        return RecordingsPayload(**raw_data)
    except Exception as e:
        raise ValueError(f"Invalid payload structure: {e}")


def _perform_transcription(s3_key: str) -> str:
    logger.info(f"Rozpoczynam transkrypcję pliku: {s3_key}")

    text = _transcription_service.process_recording(key=s3_key, bucket=settings.S3_RECORDINGS_BUCKET_NAME)

    if not text:
        raise RuntimeError("Transkrypcja zwróciła pusty wynik.")

    logger.info(f"Transkrypcja gotowa. Długość: {len(text)} znaków.")
    return text


def _perform_summarization(transcript: str, class_id: str, meeting_id: str) -> None:
    logger.info("Generowanie podsumowania...")

    _summarization_service.generate_and_save_outputs(
        transcript=transcript,
        class_id=class_id,
        meeting_id=meeting_id
    )
    logger.info("Podsumowanie wygenerowane i zapisane w S3.")


def _trigger_cleanup(event: RecordingsPayload) -> None:
    files = [event.merged] + event.voices
    logger.info(f"Zlecanie usunięcia {len(files)} plików.")

    delete_audio_files_task.delay({"file_paths": files})


def _send_notification(class_id: str, meeting_id: str, status: ProcessingStatus) -> None:
    logger.info(f"Wysyłanie powiadomienia: {status}")
    _notification_publisher.publish_resources_generated(class_id, meeting_id, status)
