package com.github.ecboot.service.admin.impl;

import com.github.ecboot.entity.AdminUser;
import com.github.ecboot.request.console.AuthLoginRequest;
import com.github.ecboot.service.admin.AdminUserService;
import com.github.ecboot.support.JsonResponse;
import org.springframework.stereotype.Service;

@Service
public class AdminUserServiceImpl implements AdminUserService {

    public JsonResponse<AdminUser> login(AuthLoginRequest authLoginRequest) {
        AdminUser adminUser = new AdminUser();
        return JsonResponse.succeed(adminUser);
    }
}
