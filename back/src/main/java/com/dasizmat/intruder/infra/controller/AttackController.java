package com.dasizmat.intruder.infra.controller;

import com.dasizmat.intruder.infra.controller.dto.AttackInputDTO;
import com.dasizmat.intruder.infra.controller.dto.AttackOutputDTO;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.net.URI;
import java.net.http.*;
import java.util.List;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;
import java.util.stream.Collectors;
import java.util.stream.Stream;

@RestController()
@RequestMapping("/api/")
public class AttackController {

    private final HttpClient client = HttpClient.newBuilder()
            .executor(Executors.newVirtualThreadPerTaskExecutor())
            .build();

    @PostMapping("/attack")
    public ResponseEntity<List<>> attack(@RequestBody AttackInputDTO request) {
        List<Future<AttackOutputDTO>> futures = Stream.of(1,2,3).map(word -> {
            String finalRequest = "".replace("{{payload}}", word);
            HttpRequest httpRequest = HttpRequest.newBuilder()
                    .uri(URI.create(finalRequest))
                    .GET()
                    .build();

            return client.sendAsync(httpRequest, HttpResponse.BodyHandlers.ofString())
                    .thenApply(resp -> new AttackOutputDTO(word, resp.statusCode(), resp.body()));
        }).map(CompletableFuture::join).collect(Collectors.toList());

        List<AttackOutputDTO> responses = futures.stream().collect(Collectors.toList());
        return ResponseEntity.ok(responses);
    }
}
