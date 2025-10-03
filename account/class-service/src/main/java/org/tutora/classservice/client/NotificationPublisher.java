package org.tutora.classservice.client;

import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.tutora.classservice.event.ClassInvitationCreatedEvent;
import org.tutora.classservice.event.EventWrapper;

@Component
@RequiredArgsConstructor
@Slf4j
public class NotificationPublisher {

    private final RabbitTemplate rabbitTemplate;

    @Value("${spring.rabbitmq.template.exchange}")
    private String exchange;

    public void sendClassInvitation(ClassInvitationCreatedEvent invitation) {
        String routingKey = ClassInvitationCreatedEvent.class.getSimpleName();

        EventWrapper<ClassInvitationCreatedEvent> event = new EventWrapper<>(
                routingKey,
                invitation
        );

        rabbitTemplate.convertAndSend(exchange, routingKey, event);
        log.info("Event published to exchange='{}', routingKey='{}', payload={}",
                exchange, routingKey, event);
    }
}
