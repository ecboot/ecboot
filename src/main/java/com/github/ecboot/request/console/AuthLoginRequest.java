package com.github.ecboot.request.console;

import lombok.Data;

import javax.validation.constraints.NotBlank;

@Data
public class AuthLoginRequest {

    @NotBlank
    private String username;

    @NotBlank
    private String password;

    @NotBlank
    private String captcha;
}
