package org.tutora.userservice.client;

import lombok.RequiredArgsConstructor;
import org.springframework.amqp.rabbit.core.RabbitTemplate;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.tutora.userservice.event.UserDetailsErrorEvent;

@Component
@RequiredArgsConstructor
public class UserDetailsErrorPublisher {

    private final RabbitTemplate rabbitTemplate;

    @Value("${spring.rabbitmq.template.exchange}")
    private String userExchange;

    public void sendErrorEvent(UserDetailsErrorEvent event) {
        rabbitTemplate.convertAndSend(userExchange, "", event);
    }
}
