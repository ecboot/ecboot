package com.github.ecboot.service.auth;

import com.github.ecboot.model.UserModel;


public interface UserAuthService {

    public UserModel login(String username, String password) {
        // userRepository.login(username, password);

        UserModel userModel = new UserModel();
        userModel.setId(1L);
        userModel.setUsername(username + "___" + password);
        userModel.setAvatar("avatar img");

        return userModel;
    }
}
