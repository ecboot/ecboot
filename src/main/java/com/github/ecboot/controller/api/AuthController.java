package com.github.ecboot.controller.api;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("api.AuthController")
@RequestMapping("api")
public class AuthController {

    @PostMapping("login")
    public String login() {
        return "login api";
    }

    @PostMapping("register")
    public String register() {
        return "register api";
    }

    @PostMapping("forgot")
    public String forgot() {
        return "forgot api";
    }

    @PostMapping("reset")
    public String reset() {
        return "reset api";
    }
}