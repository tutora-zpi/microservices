package org.tutora.classservice.config;

import org.springframework.amqp.core.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitConfig {

    @Value("${spring.rabbitmq.template.exchange}")
    private String exchangeName;

    @Value("${spring.rabbitmq.template.routing-key}")
    private String routingKey;

    @Value("${spring.rabbitmq.template.default-receive-queue}")
    private String queueName;

    @Bean
    public Exchange notificationExchange() {
        return ExchangeBuilder.topicExchange(exchangeName).durable(true).build();
    }

    @Bean
    public Queue invitationQueue() {
        return QueueBuilder.durable(routingKey).build();
    }

    @Bean
    public Binding binding(Queue invitationQueue, Exchange notificationExchange) {
        return BindingBuilder
                .bind(invitationQueue)
                .to(notificationExchange)
                .with(queueName)
                .noargs();
    }
}
