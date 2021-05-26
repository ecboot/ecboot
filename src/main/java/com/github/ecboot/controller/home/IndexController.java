package com.github.ecboot.controller.home;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

@Controller("home.IndexController")
public class IndexController {

    @GetMapping("user")
    public String index() {
        return "home/index";
    }
}
