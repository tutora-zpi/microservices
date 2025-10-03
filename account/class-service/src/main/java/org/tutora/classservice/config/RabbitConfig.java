package org.tutora.classservice.config;

import org.springframework.amqp.core.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitConfig {

    @Value("${spring.rabbitmq.template.exchange}")
    private String exchangeName;

    @Bean
    public Exchange notificationExchange() {
        return ExchangeBuilder.fanoutExchange(exchangeName).durable(true).build();
    }

    @Bean
    public Queue classInvitationQueue() {
        return QueueBuilder.durable("ClassInvitationCreatedEvent").build();
    }

    @Bean
    public Binding binding(Queue classInvitationQueue, Exchange notificationExchange) {
        return BindingBuilder
                .bind(classInvitationQueue)
                .to(notificationExchange)
                .with("ClassInvitationCreatedEvent")
                .noargs();
    }
}
