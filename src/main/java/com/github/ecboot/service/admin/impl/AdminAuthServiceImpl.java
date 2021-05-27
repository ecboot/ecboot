package com.github.ecboot.service.admin.impl;

import com.github.ecboot.entity.AdminUser;
import com.github.ecboot.repository.AdminUserRepository;
import com.github.ecboot.request.console.AuthLoginRequest;
import com.github.ecboot.service.admin.AdminAuthService;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Service
public class AdminAuthServiceImpl implements AdminAuthService {

    private final AdminUserRepository adminUserRepository;

    public AdminAuthServiceImpl(AdminUserRepository adminUserRepository) {
        this.adminUserRepository = adminUserRepository;
    }

    public Boolean login(AuthLoginRequest authLoginRequest) {
        Optional<AdminUser> adminUser = adminUserRepository.findById(1L);

        return adminUser.isPresent();
    }
}
