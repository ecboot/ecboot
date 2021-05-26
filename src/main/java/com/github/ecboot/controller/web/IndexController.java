package com.github.ecboot.controller.web;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

@Controller("web.IndexController")
public class IndexController {

    @GetMapping("/")
    public String index() {
        return "index";
    }
}
