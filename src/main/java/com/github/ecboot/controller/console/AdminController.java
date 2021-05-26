package com.github.ecboot.controller.console;

import com.github.ecboot.service.admin.AdminLogService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.ResponseBody;

@Controller("admin.admin")
@RequestMapping("admin/admin")
public class AdminController {

    @Autowired
    private AdminLogService adminLogService;

    @GetMapping("test")
    @ResponseBody
    public String show() {
        return "user show";
    }
}
