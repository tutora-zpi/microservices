import logging
from typing import List

from kombu import Connection, Exchange
from app.core.config import settings
from app.schemas.events import (
    NotificationEvent,
    ResourcesGeneratedData,
    ProcessingStatus,
)

logger = logging.getLogger(__name__)


class NotificationPublisher:
    def __init__(self):
        self.broker_url = settings.RABBITMQ_BROKER_URL
        self.exchange_name = settings.RESOURCES_EXCHANGE_NAME

        self.exchange = Exchange(self.exchange_name, type="fanout", durable=True)

        logger.info(
            f"NotificationPublisher zainicjalizowany (Exchange: {self.exchange_name})"
        )

    def publish_resources_generated(
        self,
        class_id: str,
        meeting_id: str,
        member_ids: List[str],
        status: ProcessingStatus.SUCCESS,
    ):
        try:
            event_data = ResourcesGeneratedData(
                class_id=class_id,
                meeting_id=meeting_id,
                status=status,
                member_ids=member_ids,
            )

            notification = NotificationEvent(
                pattern="ResourcesGeneratedEvent",
                data=event_data,
            )

            with Connection(self.broker_url) as conn:
                producer = conn.Producer(serializer="json")
                dict_payload = notification.model_dump(by_alias=True)

                producer.publish(
                    dict_payload,
                    exchange=self.exchange,
                    routing_key="",
                    declare=[self.exchange],
                    retry=True,
                )

            logger.info(f"Wysłano powiadomienie dla meetingId: {meeting_id}")

        except Exception as e:
            logger.error(
                f"Nie udało się wysłać powiadomienia do RabbitMQ: {e}", exc_info=True
            )
