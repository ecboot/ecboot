package com.github.ecboot.controller.home;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ResponseBody;

import javax.servlet.http.HttpSession;

@Controller("home.IndexController")
public class IndexController {

    @GetMapping("user")
    public String index() {
        return "home/index";
    }

    @GetMapping("/session")
    @ResponseBody
    public String session(HttpSession session) {
        session.setAttribute("auth", 1);

        return session.getId() + ": " + session.getAttribute("auth");
    }
}
