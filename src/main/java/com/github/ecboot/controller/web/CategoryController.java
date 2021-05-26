package com.github.ecboot.controller.web;

import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller
@Slf4j
public class CategoryController {

    @GetMapping("/category/{id}")
    @ResponseBody
    public String index(@PathVariable("id") Long id) {
        return "category page" + id;
    }
}
