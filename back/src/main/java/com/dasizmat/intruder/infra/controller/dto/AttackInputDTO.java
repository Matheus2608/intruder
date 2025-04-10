package com.dasizmat.intruder.infra.controller.dto;

import com.dasizmat.intruder.domain.TypeOfAttack;

import java.util.List;

public record AttackInputDTO(
    TypeOfAttack typeOfAttack,
    String RequestData,
    List<List<String>> payload
) { }
