package org.tutora.userservice.config;

import org.springframework.amqp.core.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration
public class RabbitConfig {

    @Value("${spring.rabbitmq.template.exchange}")
    private String exchangeName;

    @Bean
    public FanoutExchange notificationsExchange() {
        return ExchangeBuilder
                .fanoutExchange(exchangeName)
                .durable(true)
                .build();
    }

    @Bean
    public Queue userDetailsRequestedQueue() {
        return QueueBuilder
                .durable("user-details-requested-queue")
                .build();
    }

    @Bean
    public Binding userDetailsRequestedBinding(FanoutExchange notificationsExchange,
                                               Queue userDetailsRequestedQueue) {
        return BindingBuilder
                .bind(userDetailsRequestedQueue)
                .to(notificationsExchange);
    }
}
