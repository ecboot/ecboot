package com.github.ecboot.interceptor;

import cn.hutool.core.util.URLUtil;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.HttpSession;

public class AuthenticateInterceptor implements HandlerInterceptor {

    protected String guard;

    public AuthenticateInterceptor(String guard) {
        this.guard = guard;
    }

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception {
        HttpSession session = request.getSession();

        if (session.getAttribute(this.guard + "_auth_id") == null) {
            String redirectUrl = this.guard.equals("admin") ? this.guard + "/login" : "login";
            String callback = URLUtil.encodeQuery(request.getRequestURL().toString());
            response.sendRedirect("/" + redirectUrl + "?callback=" + callback);

            return false;
        }

        return true;
    }

    @Override
    public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView) throws Exception {
        HandlerInterceptor.super.postHandle(request, response, handler, modelAndView);
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex) throws Exception {
        HandlerInterceptor.super.afterCompletion(request, response, handler, ex);
    }
}
