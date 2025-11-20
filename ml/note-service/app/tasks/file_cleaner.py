from celery import shared_task
from celery.utils.log import get_task_logger
from app.schemas.events import DeleteAudioPayload
from app.services.storage import StorageS3

logger = get_task_logger(__name__)


@shared_task(name='app.tasks.file_cleaner.delete_audio_files')
def delete_audio_files_task(payload: dict):
    logger.info("Rozpoczęto zadanie czyszczenia plików (Cleanup Task).")

    try:
        data = DeleteAudioPayload(**payload)
        storage = StorageS3()

        for path in data.file_paths:
            logger.debug(f"Przetwarzanie ścieżki do usunięcia: {path}")
            storage.delete_audio(path)

        logger.info(f"Zakończono cleanup. Przetworzono {len(data.file_paths)} ścieżek.")

    except Exception as e:
        logger.error(f"Błąd w zadaniu cleanup: {e}", exc_info=True)