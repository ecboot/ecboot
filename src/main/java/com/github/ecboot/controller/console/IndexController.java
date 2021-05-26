package com.github.ecboot.controller.console;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

@Controller("console.IndexController")
public class IndexController {

    @GetMapping("admin")
    public String index() {
        return "console/index";
    }
}
