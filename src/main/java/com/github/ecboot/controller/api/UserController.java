package com.github.ecboot.controller.api;

import com.github.ecboot.support.JsonResponse;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController("api.UserController")
@RequestMapping("api/user")
public class UserController {

    @GetMapping("profile")
    public JsonResponse<String> profile() {
        return JsonResponse.succeed("user profile");
    }
}