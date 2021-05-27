package com.github.ecboot.controller.api;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("api.OAuthController")
@RequestMapping("api/oauth")
public class OAuthController {

    @GetMapping("authorize")
    public String authorize() {
        return "authorize url";
    }

    @GetMapping("callback")
    public String callback() {
        return "authorize url";
    }
}