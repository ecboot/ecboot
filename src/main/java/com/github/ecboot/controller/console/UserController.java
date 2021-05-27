package com.github.ecboot.controller.console;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;

import javax.servlet.http.HttpSession;

@Controller("console.UserController")
@RequestMapping("/admin/user")
public class UserController {

    @GetMapping("/show")
    @ResponseBody
    public String show(HttpSession httpSession) {

        return "user show: " + httpSession.getId();
    }
}