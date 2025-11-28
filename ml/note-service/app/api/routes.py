from fastapi import APIRouter, HTTPException, Path, Depends
from fastapi.responses import RedirectResponse

from app.schemas.downloads import FileType, PresignedUrlResponse, FileListResponse
from app.services.downloads_service import DownloadService
from app.api.deps import get_download_service, verify_token

router = APIRouter(dependencies=[Depends(verify_token)])


@router.get("/list/{class_id}", response_model=FileListResponse)
async def list_available_files_for_class(
        class_id: str = Path(..., description="ID klasy"),
        download_service: DownloadService = Depends(get_download_service)
):
    files = download_service.list_all_files_for_class(class_id)

    return FileListResponse(files=files)


@router.get("/presigned/{class_id}/{file_type}/{file_id}", response_model=PresignedUrlResponse)
async def get_presigned_url(
    class_id: str = Path(..., description="ID klasy"),
    file_type: FileType = Path(..., description="Typ pliku do pobrania"),
    file_id: str = Path(..., description="ID pliku (np. UUID spotkania)"),
    download_service: DownloadService = Depends(get_download_service)
):
    """
    Generuje JSON z linkiem do pobrania.
    URL: /api/presigned/{class_id}/{file_type}/{file_id}
    """
    url = download_service.generate_presigned_link(
        file_type=file_type,
        file_id=file_id,
        class_id=class_id,
        expiration=3600
    )

    if not url:
        raise HTTPException(status_code=404, detail="Nie udało się wygenerować linku.")

    return PresignedUrlResponse(
        url=url,
        file_type=file_type,
        file_id=file_id
    )


@router.get("/download/{class_id}/{file_type}/{file_id}")
async def download_file_direct(
        class_id: str = Path(..., description="ID klasy"),
        file_type: FileType = Path(..., description="Typ pliku do pobrania"),
        file_id: str = Path(..., description="ID pliku (np. UUID spotkania)"),
        download_service: DownloadService = Depends(get_download_service),
):
    """
    Bezpośrednie przekierowanie do pliku na S3.
    Wymaga nagłówka 'Authorization: Bearer <token>'.
    """
    url = download_service.generate_presigned_link(
        class_id=class_id,
        file_type=file_type,
        file_id=file_id,
        expiration=60
    )

    if not url:
        raise HTTPException(status_code=404, detail="Plik nie został znaleziony lub błąd S3.")

    return RedirectResponse(url=url, status_code=307)
