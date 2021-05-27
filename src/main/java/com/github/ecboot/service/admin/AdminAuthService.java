package com.github.ecboot.service.admin;

import com.github.ecboot.request.console.AuthLoginRequest;

public interface AdminAuthService {

    Boolean login(AuthLoginRequest authLoginRequest);
}
