package com.github.ecboot.controller.web;

import com.github.ecboot.entity.Article;
import com.github.ecboot.model.ArticleModel;
import com.github.ecboot.service.article.ArticleService;
import org.springframework.beans.BeanUtils;
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
        List<ArticleModel> all = articleService.getAll();
        model.addAttribute("articles", all);

        System.out.println(model);
        return "index";
    }

    @GetMapping("/articles/create")
    public String create() {
        return "create";
    }

    @PostMapping("/articles")
    @ResponseBody
    public String store(@RequestBody ArticleModel articleModel) {
        System.out.println(articleModel.toString());
        articleService.save(articleModel);
        return "created";
    }

    @GetMapping("/articles/{article}")
    public ArticleModel show(@PathVariable("article") Long id) {
        Article article = articleService.show(id);

        ArticleModel articleModel = new ArticleModel();
        BeanUtils.copyProperties(article, articleModel);

        return articleModel;
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
