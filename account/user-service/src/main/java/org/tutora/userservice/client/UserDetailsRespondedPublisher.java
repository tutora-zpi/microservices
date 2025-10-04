package org.tutora.userservice.client;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.tutora.userservice.entity.User;
import org.tutora.userservice.event.EventWrapper;
import org.tutora.userservice.event.UserDetailsRespondedEvent;

@Slf4j
@Component
@RequiredArgsConstructor
public class UserDetailsRespondedPublisher {

    private final RabbitTemplate rabbitTemplate;

    @Value("${spring.rabbitmq.template.exchange}")
    private String userExchange;

    public void sendUserDetailsRespondedEvent(
            String notificationId,
            User sender,
            User receiver
    ) {
        var senderDetails = new UserDetailsRespondedEvent.UserDetails(
                sender.getId().toString(),
                sender.getName(),
                sender.getSurname(),
                sender.getRoles().stream()
                        .map(role -> role.getName().toString())
                        .toList()
        );
        var receiverDetails = new UserDetailsRespondedEvent.UserDetails(
                receiver.getId().toString(),
                receiver.getName(),
                receiver.getSurname(),
                receiver.getRoles().stream()
                        .map(role -> role.getName().toString())
                        .toList()
        );

        var wrapper = new EventWrapper<>(
                UserDetailsRespondedEvent.class.getSimpleName(),
                new UserDetailsRespondedEvent(notificationId, senderDetails, receiverDetails)
        );
        rabbitTemplate.convertAndSend(userExchange, "", wrapper);
        log.info("Published UserDetailsRespondedEvent for notificationId={} -> {}", notificationId, wrapper);
    }
}
