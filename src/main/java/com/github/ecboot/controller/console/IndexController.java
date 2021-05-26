package com.github.ecboot.controller.console;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;

@Controller("console.IndexController")
public class IndexController {

    @GetMapping("/admin")
    public String index(HttpSession session) {
        HashMap<String, String> map = new HashMap<>();
        UUID uid = (UUID) session.getAttribute("uid");

        if (uid == null) {
            System.out.println("not login");
            uid = UUID.randomUUID();
            session.setAttribute("uid", uid);
        } else {
            System.out.println("uid: " + uid);
        }

        return "admin/index";
    }
}
