package com.github.ecboot.controller.web;

import cn.hutool.captcha.CaptchaUtil;
import cn.hutool.captcha.LineCaptcha;
import cn.hutool.captcha.generator.RandomGenerator;
import com.github.ecboot.constant.GlobalConst;
import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.ResponseBody;

import javax.servlet.ServletOutputStream;
import javax.servlet.http.HttpServletResponse;
import javax.servlet.http.HttpSession;
import java.io.IOException;

@Controller
public class CaptchaController {

    @GetMapping("/captcha")
    @ResponseBody
    public void index(HttpServletResponse response, HttpSession session) throws IOException {
        RandomGenerator randomGenerator = new RandomGenerator("023456789", 4);

        LineCaptcha captcha = CaptchaUtil.createLineCaptcha(90, 38);
        captcha.setGenerator(randomGenerator);
        captcha.createCode();

        session.setAttribute(GlobalConst.CAPTCHA, captcha.getCode());

        ServletOutputStream outputStream = response.getOutputStream();
        captcha.write(outputStream);
        outputStream.flush();
        outputStream.close();
    }
}
