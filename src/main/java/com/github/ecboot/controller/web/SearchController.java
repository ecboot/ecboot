package com.github.ecboot.controller.web;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
public class SearchController {

    @GetMapping("/search")
    @ResponseBody
    public String index() {
        return "search page";
    }
}
