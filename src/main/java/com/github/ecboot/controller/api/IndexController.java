package com.github.ecboot.controller.api;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("api.IndexController")
public class IndexController {

    @GetMapping("api")
    public String index() {
        return "API JSON";
    }
}
