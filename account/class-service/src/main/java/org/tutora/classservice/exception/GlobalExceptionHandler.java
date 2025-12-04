package org.tutora.classservice.exception;

import org.tutora.classservice.dto.ErrorDetailsDto;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;

import java.time.LocalDateTime;

@ControllerAdvice
public class GlobalExceptionHandler {

    private ResponseEntity<ErrorDetailsDto> buildErrorResponse(
            Exception exception,
            WebRequest webRequest,
            HttpStatus status
    ) {
        ErrorDetailsDto errorDetails = new ErrorDetailsDto(
                LocalDateTime.now(),
                exception.getMessage(),
                webRequest.getDescription(false).replace("uri=", ""),
                status.value()
        );

        return new ResponseEntity<>(errorDetails, status);
    }

    @ExceptionHandler(ResourceNotFoundException.class)
    public ResponseEntity<ErrorDetailsDto> handleResourceNotFoundException(
            ResourceNotFoundException exception,
            WebRequest webRequest) {
        return buildErrorResponse(exception, webRequest, HttpStatus.NOT_FOUND);
    }

    @ExceptionHandler(UserAlreadyInClassException.class)
    public ResponseEntity<ErrorDetailsDto> handleUserAlreadyInClassException(
            UserAlreadyInClassException exception,
            WebRequest webRequest) {
        return buildErrorResponse(exception, webRequest, HttpStatus.CONFLICT);
    }

    @ExceptionHandler(UserAlreadyInvitedException.class)
    public ResponseEntity<ErrorDetailsDto> handleUserAlreadyInvitedException(
            UserAlreadyInvitedException exception,
            WebRequest webRequest) {
        return buildErrorResponse(exception, webRequest, HttpStatus.CONFLICT);
    }

    @ExceptionHandler(UserRejectedInvitationException.class)
    public ResponseEntity<ErrorDetailsDto> handleUserRejectedInvitationException(
            UserRejectedInvitationException exception,
            WebRequest webRequest) {
        return buildErrorResponse(exception, webRequest, HttpStatus.BAD_REQUEST);
    }

    @ExceptionHandler(UnauthorizedActionException.class)
    public ResponseEntity<ErrorDetailsDto> handleUnauthorizedActionException(
            UnauthorizedActionException exception,
            WebRequest webRequest) {
        return buildErrorResponse(exception, webRequest, HttpStatus.FORBIDDEN);
    }

    @ExceptionHandler(IllegalArgumentException.class)
    public ResponseEntity<ErrorDetailsDto> handleIllegalArgumentException(
            UnauthorizedActionException exception,
            WebRequest webRequest) {
        return buildErrorResponse(exception, webRequest, HttpStatus.BAD_REQUEST);
    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorDetailsDto> handleGlobalException(
            Exception exception,
            WebRequest webRequest) {
        return buildErrorResponse(
                new Exception("An unexpected internal server error occurred: " + exception.getMessage()),
                webRequest,
                HttpStatus.INTERNAL_SERVER_ERROR
        );
    }
}
