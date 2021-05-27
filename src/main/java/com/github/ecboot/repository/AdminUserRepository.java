package com.github.ecboot.repository;

import com.github.ecboot.entity.AdminUser;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

public interface AdminUserRepository extends JpaRepository<AdminUser, Long> {

    AdminUser getByUserName(String username);
}
