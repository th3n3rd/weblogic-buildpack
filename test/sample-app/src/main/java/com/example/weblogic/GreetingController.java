package com.example.weblogic;

import lombok.Value;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

@RestController
@ResponseBody
public class GreetingController {

    @GetMapping("/")
    public GreetResponse greet() {
        return new GreetResponse("Hello World!");
    }

    @Value
    static class GreetResponse {
        String message;
    }
}
