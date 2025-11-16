from celery import Celery
from app.core.config import settings

broker_url = settings.RABBITMQ_BROKER_URL
backend_url = "rpc://"

celery_app = Celery(
    "tasks",
    broker=broker_url,
    backend=backend_url,
    include=['app.tasks.meeting_processor']
)
