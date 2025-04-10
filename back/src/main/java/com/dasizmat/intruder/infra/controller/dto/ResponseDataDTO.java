package com.dasizmat.intruder.infra.controller.dto;

public record ResponseDataDTO(
        Integer requestId,
        String payload,
        Integer statusCode,
        Integer elapsedTimeInMilliseconds,
        Boolean didError,
        Integer bytesLenght,
        String request,
        String response
) { }
