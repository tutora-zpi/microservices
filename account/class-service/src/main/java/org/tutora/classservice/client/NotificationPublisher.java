package org.tutora.classservice.client;

import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.stereotype.Component;
import org.tutora.classservice.dto.ClassInvitationEvent;
import org.tutora.classservice.dto.Event;

@Component
@RequiredArgsConstructor
public class NotificationPublisher {

    private final RabbitTemplate rabbitTemplate;

    public void sendClassInvitation(ClassInvitationEvent invitation) {
        Event<ClassInvitationEvent> event = new Event<>(
                ClassInvitationEvent.class.getSimpleName(),
                invitation
        );

        rabbitTemplate.convertAndSend(event);
    }
}
