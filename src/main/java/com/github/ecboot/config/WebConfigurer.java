package com.github.ecboot.config;

import com.github.ecboot.interceptor.AuthenticateInterceptor;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

import java.util.Arrays;
import java.util.List;

@Configuration
public class WebConfigurer implements WebMvcConfigurer {

    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        // 多个拦截器时 以此添加 执行顺序按添加顺序
        List<String> guards = Arrays.asList("admin", "user");

        // 注册拦截器 拦截规则
        for (String guard : guards) {
            registry.addInterceptor(getHandlerInterceptor(guard))
                    .addPathPatterns("/" + guard + "/**", "/" + guard)
                    .excludePathPatterns("/" + guard + "/login", "/" + guard + "/forgot", "/" + guard + "/reset");
        }
    }

    public static HandlerInterceptor getHandlerInterceptor(String guard) {
        return new AuthenticateInterceptor(guard);
    }
}
