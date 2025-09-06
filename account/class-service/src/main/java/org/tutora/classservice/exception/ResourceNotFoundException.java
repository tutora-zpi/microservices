package org.tutora.classservice.exception;

import java.util.Map;
import java.util.stream.Collectors;

public class ResourceNotFoundException extends RuntimeException {

  public ResourceNotFoundException(String message) {
    super(message);
  }

  public ResourceNotFoundException(String resourceName, String fieldName, Object fieldValue) {
    super(String.format("%s not found with %s: '%s'", resourceName, fieldName, fieldValue));
  }

  public ResourceNotFoundException(String resourceName, Map<String, Object> fields) {
    super(String.format("%s not found with %s", resourceName,
            fields.entrySet().stream()
                    .map(e -> e.getKey() + ": '" + e.getValue() + "'")
                    .collect(Collectors.joining(", "))));
  }
}