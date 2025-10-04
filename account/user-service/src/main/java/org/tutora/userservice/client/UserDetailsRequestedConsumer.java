package org.tutora.userservice.client;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.springframework.amqp.rabbit.annotation.RabbitListener;
import org.springframework.stereotype.Component;
import org.tutora.userservice.event.EventWrapper;
import org.tutora.userservice.event.UserDetailsErrorEvent;
import org.tutora.userservice.event.UserDetailsRequestedEvent;
import org.tutora.userservice.exception.ResourceNotFoundException;
import org.tutora.userservice.service.contract.UserService;

import java.util.UUID;

@Slf4j
@Component
@RequiredArgsConstructor
public class UserDetailsRequestedConsumer {

    private final UserDetailsRespondedPublisher publisher;
    private final UserDetailsErrorPublisher errorPublisher;
    private final UserService userService;

    @RabbitListener(queues = "user-details-requested-queue")
    public void handleUserDetailsRequest(EventWrapper<UserDetailsRequestedEvent> wrapper) {
        if (!UserDetailsRequestedEvent.class.getSimpleName().equals(wrapper.pattern())) {
            log.debug("Ignoring event type: {}", wrapper.pattern());
            return;
        }

        var event = wrapper.data();
        log.info("Received UserDetailsRequestedEvent: {}", event);
        if (event.notificationId() == null || event.notificationId().isEmpty()) {
            log.warn("Received event with null notificationId, skipping processing");
            return;
        }

        try {
            var sender = userService.findById(UUID.fromString(event.senderId()));
            var receiver = userService.findById(UUID.fromString(event.receiverId()));

            publisher.sendUserDetailsRespondedEvent(
                    event.notificationId(),
                    sender,
                    receiver
            );
        } catch (ResourceNotFoundException e) {
            log.warn("User not found: {}", e.getMessage());

            errorPublisher.sendErrorEvent(new UserDetailsErrorEvent(
                    event.notificationId(),
                    e.getMessage()
            ));
        }
    }
}
