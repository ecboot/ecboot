package com.github.ecboot.service.user;

import com.github.ecboot.model.UserModel;

public interface UserAuthService {

    Boolean login(String username, String password);
}
