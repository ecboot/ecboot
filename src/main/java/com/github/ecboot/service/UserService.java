package com.github.ecboot.service;

public interface UserService {
    /**
     * 用户登录
     *
     * @param username 用户名/手机号/电子邮箱
     * @param password 登录密码
     * @return Boolean
     */
    Boolean login(String username, String password);
}
