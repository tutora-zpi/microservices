from celery import Celery
from celery import bootsteps
from kombu import Consumer, Exchange, Queue
from app.core.config import settings
import logging

broker_url = settings.RABBITMQ_BROKER_URL
backend_url = "rpc://"

celery_app = Celery(
    "tasks",
    broker=broker_url,
    backend=backend_url,
    include=['app.tasks.meeting_processor']
)

meeting_exchange = Exchange('meeting', type='fanout')

meeting_queue = Queue(
    name='celery-meeting-listener',
    exchange=meeting_exchange,
    routing_key='#',
    durable=False,
    auto_delete=False
)


class EventConsumer(bootsteps.ConsumerStep):
    def get_consumers(self, channel):
        return [Consumer(channel,
                         queues=[meeting_queue],
                         callbacks=[self.handle_message],
                         accept=['json'])]

    def handle_message(self, body, message):
        event_name = body.get('pattern')
        payload = body.get('data')

        logging.info(f"Otrzymano event: {event_name}")

        if event_name == 'RecordingsUploaded':
            from app.tasks.meeting_processor import process_audio_file_task
            logging.info(f"Payload: {payload}")

            process_audio_file_task.delay(payload)

        else:
            logging.warning(f"Nieznany event: {event_name}")

        message.ack()


celery_app.steps['consumer'].add(EventConsumer)
