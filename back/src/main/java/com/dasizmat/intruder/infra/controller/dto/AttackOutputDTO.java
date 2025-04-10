package com.dasizmat.intruder.infra.controller.dto;

import java.util.List;

public record AttackOutputDTO(
        List<ResponseDataDTO> responses,
        String url,
        Integer totalElapsedTimeInMiliseconds
) { }
