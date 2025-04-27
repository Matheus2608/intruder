package com.dasizmat.intruder.infra.controller;

import com.dasizmat.intruder.infra.controller.dto.AttackInputDTO;
import com.dasizmat.intruder.infra.controller.dto.AttackOutputDTO;
import com.dasizmat.intruder.infra.controller.dto.ResponseDataDTO;
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

    @GetMapping("/fake")
    @CrossOrigin(origins = "http://localhost:4200")
    public ResponseEntity<AttackOutputDTO> fakeApi() {
        return ResponseEntity.ok(createFakeOutput());
    }

    @PostMapping("/fake")
    @CrossOrigin(origins = "http://localhost:4200")
    public ResponseEntity<AttackOutputDTO> postFakeApi(@RequestBody AttackInputDTO body) {
        System.out.println(body);
        return ResponseEntity.ok(createFakeOutput());
    }

    AttackOutputDTO createFakeOutput() {
        return new AttackOutputDTO(
                getFakeResposes(),
                "http://localhost:8080/api/attack",
                1234
        );
    }

    List<ResponseDataDTO> getFakeResposes() {
        return List.of(
                new ResponseDataDTO(1, "Payload 1", 200, 123, false, 512, "GET /api/data", "{\"data\":\"value1\"}"),
                new ResponseDataDTO(2, "Payload 2", 404, 256, true, 0, "GET /api/unknown", "{\"error\":\"Not Found\"}"),
                new ResponseDataDTO(3, "Payload 3", 500, 500, true, 1024, "POST /api/create", "{\"error\":\"Internal Server Error\"}"),
                new ResponseDataDTO(4, "Payload 4", 201, 150, false, 768, "POST /api/create", "{\"data\":\"created\"}"),
                new ResponseDataDTO(5, "Payload 5", 400, 320, true, 256, "POST /api/invalid", "{\"error\":\"Bad Request\"}"),
                new ResponseDataDTO(6, "Payload 6", 200, 87, false, 1024, "GET /api/status", "{\"status\":\"ok\"}"),
                new ResponseDataDTO(7, "Payload 7", 403, 210, true, 128, "GET /api/forbidden", "{\"error\":\"Forbidden\"}"),
                new ResponseDataDTO(8, "Payload 8", 204, 95, false, 0, "DELETE /api/resource", "{}"),
                new ResponseDataDTO(9, "Payload 9", 202, 180, false, 300, "POST /api/async", "{\"message\":\"Accepted\"}"),
                new ResponseDataDTO(10, "Payload 10", 503, 600, true, 0, "GET /api/down", "{\"error\":\"Service Unavailable\"}")
        );
    }

//    @PostMapping("/attack")
//    public ResponseEntity<List<>> attack(@RequestBody AttackInputDTO request) {
//        List<Future<AttackOutputDTO>> futures = Stream.of(1,2,3).map(word -> {
//            String finalRequest = "".replace("{{payload}}", word);
//            HttpRequest httpRequest = HttpRequest.newBuilder()
//                    .uri(URI.create(finalRequest))
//                    .GET()
//                    .build();
//
//            return client.sendAsync(httpRequest, HttpResponse.BodyHandlers.ofString())
//                    .thenApply(resp -> new AttackOutputDTO(word, resp.statusCode(), resp.body()));
//        }).map(CompletableFuture::join).collect(Collectors.toList());
//
//        List<AttackOutputDTO> responses = futures.stream().collect(Collectors.toList());
//        return ResponseEntity.ok(responses);
//    }
}
