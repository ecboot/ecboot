package com.github.ecboot.repository;

import com.github.ecboot.entity.Users;
import org.springframework.data.jpa.repository.JpaRepository;

public interface UsersRepository extends JpaRepository<Users, Long> {

    // UserModel login(String username, String password);
}
