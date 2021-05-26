package com.github.ecboot.request;

import lombok.Data;

import javax.validation.constraints.NotBlank;

@Data
public class LoginRequest {

    @NotBlank(message = "请填写登录用户名")
    private String username;

    @NotBlank(message = "请填写登录密码")
    private String password;

    @NotBlank(message = "请填写图片验证码")
    private String captcha;
}
