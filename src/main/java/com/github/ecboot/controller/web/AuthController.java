package com.github.ecboot.controller.web;

import com.github.ecboot.model.UserModel;
import com.github.ecboot.request.UserLoginRequest;
import com.github.ecboot.service.auth.UserAuthService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.util.DigestUtils;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.ResponseBody;

import javax.validation.Valid;

@Controller
public class AuthController {

    @Autowired
    UserAuthService authService;

    @GetMapping("/login")
    public String index() {
        return "auth/login";
    }

    @PostMapping("/login")
    @ResponseBody
    public UserModel login(@Valid @RequestBody UserLoginRequest userLoginRequest,
                           BindingResult bindingResult) {
        if (bindingResult.hasErrors()) {
            // return "error: " + Objects.requireNonNull(bindingResult.getFieldError()).getDefaultMessage();
        }

        userLoginRequest.setPassword(DigestUtils.md5DigestAsHex(userLoginRequest.getPassword().getBytes()));

        return authService.login(userLoginRequest.getUsername(), userLoginRequest.getPassword());
    }
}
