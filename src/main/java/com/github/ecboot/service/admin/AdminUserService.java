package com.github.ecboot.service.admin;

import com.github.ecboot.entity.AdminUser;
import com.github.ecboot.request.AdminLoginRequest;
import com.github.ecboot.utils.JsonResponse;


public interface AdminUserService {

    public JsonResponse<AdminUser> login(AdminLoginRequest adminLoginRequest) {
        AdminUser adminUser = new AdminUser();
        return JsonResponse.succeed(adminUser);
    }
}
