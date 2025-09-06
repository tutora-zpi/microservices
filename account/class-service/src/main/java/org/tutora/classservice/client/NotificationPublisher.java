package org.tutora.classservice.client;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;
import org.tutora.classservice.dto.ClassInvitationEvent;
import org.tutora.classservice.dto.Event;

import java.io.UncheckedIOException;

@Component
@RequiredArgsConstructor
public class NotificationPublisher {

    private final RabbitTemplate rabbitTemplate;
    private final ObjectMapper objectMapper;

    public void sendClassInvitation(ClassInvitationEvent invitation) {
        Event<ClassInvitationEvent> event = new Event<>(
                ClassInvitationEvent.class.getSimpleName(),
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
