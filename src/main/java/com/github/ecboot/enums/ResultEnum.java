package com.github.ecboot.enums;

import lombok.Getter;

/**
 * 请求状态枚举
 * 后续状态若太多，可以考虑使用继承等方式拓展。
 * 10x: message
 * 20x: success
 * 30x: redirect
 * 40x: request error
 * 50x: server error
 * 600: unparseable response headers
 */
@Getter
public enum ResultEnum {

    /**
     * 通用
     */
    ERROR(-1, "服务端错误"),
    SUCCESS(0, "成功"),

    CAPTCHA_ERROR(400001, "图片验证码错误"),
    FIELD_ERROR(400002, "字段错误"),
    DUPLICATE_EXCEPTION(400999, "重复异常"),

    /**
     * 用户模块
     */
    PASSWORD_ERROR(1, "密码错误"),
    USERNAME_EXIST(2, "用户名已存在"),
    USERNAME_OR_PASSWORD_ERROR(400101, "用户名或密码错误"),

    ;

    private final Integer code;

    private final String message;

    ResultEnum(Integer code, String message) {
        this.code = code;
        this.message = message;
    }
}