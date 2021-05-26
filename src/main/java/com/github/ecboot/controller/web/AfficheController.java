package com.github.ecboot.controller.web;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
public class AfficheController {

    @GetMapping("/affiche")
    @ResponseBody
    public String index() {
        return "AfficheController page";
    }
}
