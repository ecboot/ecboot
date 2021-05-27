package com.github.ecboot.controller.home;

import com.github.ecboot.constant.GlobalConst;
import com.github.ecboot.entity.AdminUser;
import com.github.ecboot.enums.ResultEnum;
import com.github.ecboot.request.console.AuthLoginRequest;
import com.github.ecboot.service.admin.AdminAuthService;
import com.github.ecboot.service.admin.AdminUserService;
import com.github.ecboot.support.JsonResponse;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpSession;
import javax.validation.Valid;

@Controller("home.AuthController")
@RequestMapping("user")
public class AuthController {

    @Autowired
    private AdminAuthService adminAuthService;

    @Autowired
    private AdminUserService adminUserService;

    @GetMapping("login")
    public String index() {
        return "home/auth/login";
    }

    @PostMapping("login")
    @ResponseBody
    public JsonResponse<AdminUser> login(@Valid @RequestBody AuthLoginRequest authLoginRequest,
                                         BindingResult bindingResult,
                                         HttpSession session) {
        // 表单验证
        if (bindingResult.hasErrors()) {
            return JsonResponse.failed(ResultEnum.FIELD_ERROR, bindingResult);
        }

        // 图片验证码校验
        if (!session.getAttribute(GlobalConst.CAPTCHA).toString().equalsIgnoreCase(authLoginRequest.getCaptcha())) {
            return JsonResponse.failed(ResultEnum.CAPTCHA_ERROR);
        }

        if (adminAuthService.login(authLoginRequest)) {
            AdminUser adminUser = adminUserService.getAdminUserByUsername(authLoginRequest.getUsername());

            if (adminUser != null) {
                session.setAttribute(GlobalConst.AUTH_ADMIN_ID, adminUser.getUserId());
                // TOTO pojo transform
                return JsonResponse.succeed(adminUser);
            }
        }

        return JsonResponse.failed(ResultEnum.USERNAME_OR_PASSWORD_ERROR);
    }

    @GetMapping("forgot")
    public String forgot() {
        return "home/auth/forgot";
    }

    @GetMapping("reset")
    public String reset() {
        return "home/auth/reset";
    }

}
