package org.tutora.classservice.client;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;
import org.tutora.classservice.event.ClassInvitationCreatedEvent;
import org.tutora.classservice.event.EventWrapper;

import java.io.UncheckedIOException;

@Component
@RequiredArgsConstructor
public class NotificationPublisher {

    private final RabbitTemplate rabbitTemplate;
    private final ObjectMapper objectMapper;

    public void sendClassInvitation(ClassInvitationCreatedEvent invitation) {
        EventWrapper<ClassInvitationCreatedEvent> event = new EventWrapper<>(
                ClassInvitationCreatedEvent.class.getSimpleName(),
                invitation
        );

        try {
            String json = objectMapper.writeValueAsString(event);
            rabbitTemplate.convertAndSend(json);
        } catch (JsonProcessingException e) {
            throw new UncheckedIOException("Error serializing event: ", e);
        }
    }
}
