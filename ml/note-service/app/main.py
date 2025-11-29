from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.routes import router as downloads_router
from app.core.config import settings

app = FastAPI(
    title="TutorAI Storage API",
    description="API do zarządzania plikami i generowania linków S3",
    version="1.0.0"
)


origins = [
    "http://localhost:3000",  # Twój frontend
    "http://localhost:8080",
    "*"                       # W developmencie można dać '*', na produkcji podaj konkretne domeny
]

if settings.BACKEND_CORS_ORIGINS:
    app.add_middleware(
        CORSMiddleware,
        allow_origins=settings.BACKEND_CORS_ORIGINS,
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )


app.include_router(downloads_router, prefix="/api", tags=["Downloads"])


@app.get("/health")
def health_check():
    return {"status": "ok", "service": "note-service"}
