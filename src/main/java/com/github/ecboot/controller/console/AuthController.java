package com.github.ecboot.controller.console;

import com.github.ecboot.constant.SystemConst;
import com.github.ecboot.entity.AdminUser;
import com.github.ecboot.enums.ResultEnum;
import com.github.ecboot.request.AdminLoginRequest;
import com.github.ecboot.service.admin.AdminUserService;
import com.github.ecboot.utils.JsonResponse;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.validation.BindingResult;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpSession;
import javax.validation.Valid;

@Controller("admin.auth")
@RequestMapping("admin")
@Slf4j
public class AuthController {

    @Autowired
    private AdminUserService adminUserService;

    @GetMapping("login")
    public String index() {
        return "admin/auth/login";
    }

    @PostMapping("login")
    @ResponseBody
    public JsonResponse<AdminUser> login(@Valid @RequestBody AdminLoginRequest adminLoginRequest,
                                         BindingResult bindingResult,
                                         HttpSession session) {
        // 表单验证
        if (bindingResult.hasErrors()) {
            return JsonResponse.failed(ResultEnum.FIELD_ERROR, bindingResult);
        }

        // 图片验证码校验
        if (!session.getAttribute(SystemConst.CAPTCHA).toString().equalsIgnoreCase(adminLoginRequest.getCaptcha())) {
            return JsonResponse.failed(ResultEnum.CAPTCHA_ERROR);
        }

        JsonResponse<AdminUser> adminUser = adminUserService.login(adminLoginRequest);

        if (adminUser.getStatus().equalsIgnoreCase("success")) {
            session.setAttribute(SystemConst.AUTH_ADMIN_ID, adminUser.getData().getUserId());
            return adminUser;
        } else {
            return JsonResponse.failed(ResultEnum.USERNAME_OR_PASSWORD_ERROR);
        }
    }
}
