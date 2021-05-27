package com.github.ecboot.service.admin;

import com.github.ecboot.entity.AdminUser;

public interface AdminUserService {

    AdminUser getAdminUserByUsername(String username);
}
