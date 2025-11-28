import json
from functools import lru_cache
from typing import Annotated, Dict, Any

import jwt
import httpx
from jwt.algorithms import RSAAlgorithm
from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from app.services.downloads_service import DownloadService

from app.core.config import Settings
from app.services.storage_s3 import StorageS3

security = HTTPBearer()
jwks_cache: Dict[str, Any] = {}


@lru_cache()
def get_settings() -> Settings:
    return Settings()


def get_storage_service(settings: Settings = Depends(get_settings)) -> StorageS3:
    return StorageS3(settings=settings)


def get_download_service(storage_service: StorageS3 = Depends(get_storage_service)) -> DownloadService:
    return DownloadService(storage_service)


async def _get_public_key(kid: str, settings: Settings) -> Any:
    """
    Pobiera klucz publiczny dla danego Key ID (kid).
    Najpierw sprawdza cache, potem woła endpoint JWKS.
    """
    if kid in jwks_cache:
        return jwks_cache[kid]

    print(f"Fetching JWKS from {settings.AUTH_JWKS_URL}...")
    async with httpx.AsyncClient() as client:
        try:
            response = await client.get(settings.AUTH_JWKS_URL)
            response.raise_for_status()
            jwks = response.json()
        except Exception as e:
            print(f"Failed to fetch JWKS: {e}")
            raise HTTPException(status_code=503, detail="Could not fetch public keys")

    found_key = None
    for key_data in jwks.get("keys", []):
        current_kid = key_data.get("kid")
        if not current_kid:
            continue

        try:
            public_key = RSAAlgorithm.from_jwk(json.dumps(key_data))
            jwks_cache[current_kid] = public_key

            if current_kid == kid:
                found_key = public_key
        except Exception as e:
            print(f"Error parsing key {current_kid}: {e}")

    return found_key


async def verify_token(
        token: Annotated[HTTPAuthorizationCredentials, Depends(security)],
        settings: Settings = Depends(get_settings)
) -> Dict[str, Any]:
    """
    Weryfikuje token JWT używając dynamicznie pobieranego JWKS.
    """
    credentials_exception = HTTPException(
        status_code=status.HTTP_401_UNAUTHORIZED,
        detail="Could not validate credentials",
        headers={"WWW-Authenticate": "Bearer"},
    )

    try:
        unverified_header = jwt.get_unverified_header(token.credentials)
        kid = unverified_header.get("kid")

        if not kid:
            print("Token missing 'kid' header")
            raise credentials_exception

        public_key = await _get_public_key(kid, settings)

        if not public_key:
            print(f"Public key not found for kid: {kid}")
            raise credentials_exception

        payload = jwt.decode(
            token.credentials,
            key=public_key,
            algorithms=[settings.AUTH_ALGORITHM],
            options={"verify_aud": False}
        )

        return payload

    except jwt.ExpiredSignatureError:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Token has expired",
            headers={"WWW-Authenticate": "Bearer"},
        )
    except jwt.PyJWTError as e:
        print(f"JWT Validation Error: {e}")
        raise credentials_exception
    except Exception as e:
        print(f"Unexpected Auth Error: {e}")
        raise credentials_exception
