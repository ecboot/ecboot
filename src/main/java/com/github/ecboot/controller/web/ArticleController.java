package com.github.ecboot.controller.web;

import com.github.ecboot.entity.Article;
import com.github.ecboot.service.article.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Controller;
import org.springframework.ui.Model;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpSession;
import java.util.List;

@Controller
@RequestMapping("articles")
public class ArticleController {

    @Autowired
    private ArticleService articleService;

    @GetMapping("/articles")
    public String index(Model model) {
        List<Article> all = articleService.getList();
        model.addAttribute("articles", all);

        System.out.println(model);
        return "index";
    }

    @GetMapping("/articles/create")
    public String create() {
        return "create";
    }

    @GetMapping("/articles/{article}")
    public Object show(@PathVariable("article") Long id) {

        return null;
    }

    @GetMapping("/articles/{article}/edit")
    public String edit(@PathVariable String article) {
        return "create";
    }

    @PutMapping("/articles/{article}")
    @ResponseBody
    public String update(@PathVariable String article) {
        return "updated";
    }

    @DeleteMapping("/articles/{article}")
    @ResponseBody
    public String destroy(@PathVariable String article) {
        return "deleted";
    }

    @GetMapping("/session")
    @ResponseBody
    public String session(HttpSession session) {
        session.setAttribute("auth", 1);

        return session.getId() + ": " + session.getAttribute("auth");
    }
}
