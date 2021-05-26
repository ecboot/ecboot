package com.github.ecboot.service.user.impl;

import com.github.ecboot.service.user.UserAuthService;
import org.springframework.stereotype.Service;

@Service
public class UserAuthServiceImpl implements UserAuthService {

    public Boolean login(String username, String password) {
        // userRepository.login(username, password);

//        UserModel userModel = new UserModel();
//        userModel.setId(1L);
//        userModel.setUsername(username + "___" + password);
//        userModel.setAvatar("avatar img");
//
        return false;
    }
}
