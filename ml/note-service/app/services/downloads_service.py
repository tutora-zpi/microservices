import os
import logging
from typing import List, Optional

from app.schemas.downloads import FileType, FileItem
from app.services.storage_s3 import StorageS3

logger = logging.getLogger(__name__)


class DownloadService:
    def __init__(self, storage_service: StorageS3):
        self.storage = storage_service

    def _resolve_s3_path(self, file_type: FileType, file_id: str, class_id: Optional[str] = None) -> tuple[str, str]:
        if file_type == FileType.recording:
            return file_id, "recordings"

        if not class_id:
            raise ValueError(f"Parametr class_id jest wymagany do pobrania zasobu typu: {file_type.value}")

        filename = f"{file_id}.md"
        folder = file_type.value
        key = f"{class_id}/{folder}/{filename}"

        return key, "notes"

    def list_all_files_for_class(self, class_id: str) -> List[FileItem]:
        folders_to_scan = [FileType.student_notes, FileType.teacher_tests]
        logger.info(folders_to_scan)
        all_files = []

        for file_type in folders_to_scan:
            folder_prefix = f"{class_id}/{file_type.value}"

            s3_objects = self.storage.list_files(folder_prefix)

            for obj in s3_objects:
                key = obj["key"]
                filename = os.path.basename(key)
                file_id = os.path.splitext(filename)[0]

                all_files.append(FileItem(
                    file_id=file_id,
                    file_type=file_type,
                    created_at=obj["last_modified"],
                    size_bytes=obj["size"]
                ))
        return all_files

    def generate_presigned_link(self, file_type: FileType, file_id: str, class_id: str, expiration: int = 3600) -> str:
        key, bucket_type = self._resolve_s3_path(file_type, file_id, class_id)

        return self.storage.generate_presigned_url(
            object_key=key,
            bucket_type=bucket_type,
            expiration=expiration
        )
