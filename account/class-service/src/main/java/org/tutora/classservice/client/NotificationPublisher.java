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
        String pattern = invitation.name();

        EventWrapper<ClassInvitationCreatedEvent> event = new EventWrapper<>(
                pattern,
                invitation
        );

        rabbitTemplate.convertAndSend(exchange, pattern, event);
        log.info("Event published to exchange='{}', pattern='{}', payload={}",
                exchange, pattern, event);
    }
}
