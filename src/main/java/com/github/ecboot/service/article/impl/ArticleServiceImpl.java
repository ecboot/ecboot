package com.github.ecboot.service.article.impl;

import com.github.ecboot.entity.Article;
import com.github.ecboot.repository.ArticleRepository;
import com.github.ecboot.service.article.ArticleService;
import org.springframework.beans.BeanUtils;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.ArrayList;
import java.util.List;

@Service
public class ArticleServiceImpl implements ArticleService {

    @Autowired
    private ArticleRepository articleRepository;

    public ArticleModel save(ArticleModel articleModel) {
        Article article = new Article();

        BeanUtils.copyProperties(articleModel, article);

        articleRepository.save(article);

        System.out.println(articleModel.toString());
        System.out.println(article.toString());
        return null;
    }

    public void deleteArticle(Long id) {

    }

    public void updateArticle(ArticleModel article) {

    }

    public Article show(Long id) {
        return articleRepository.getOne(id);
    }

    public List<ArticleModel> getAll() {
        List<ArticleModel> articleModels = new ArrayList<>();

        List<Article> all = articleRepository.findAll();
        for (Article post : all) {
            ArticleModel articleModel = new ArticleModel();
            BeanUtils.copyProperties(post, articleModel);
            articleModels.add(articleModel);
        }

        return articleModels;
    }
}
