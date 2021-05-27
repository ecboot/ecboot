package com.github.ecboot.service.admin;

import com.github.ecboot.entity.AdminUser;

public interface AdminUserService {

    Boolean createAdminUser();

    Boolean updateAdminUser();

    Boolean deleteAdminUserById(Long id);

    Boolean deleteAdminUserByUsername(String username);

    AdminUser getAdminUserById(Long id);

    AdminUser getAdminUserByUsername(String username);
}
