package com.dasizmat.intruder.infra.controller.dto;

import com.dasizmat.intruder.domain.TypeOfAttack;

import java.util.List;

public record AttackInputDTO(
    TypeOfAttack typeOfAttack,
    String requestData,
    List<List<String>> payloads
) { }
