package com.github.ecboot.service.article.impl;

import com.github.ecboot.entity.Article;
import com.github.ecboot.repository.ArticleRepository;
import com.github.ecboot.service.article.ArticleService;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;

@Service
public class ArticleServiceImpl implements ArticleService {

    @Autowired
    private ArticleRepository articleRepository;

    @Override
    public List<Article> getList() {
        // TODO date formatting
        return articleRepository.findAll();
    }

    @Override
    public Article create(Article article) {
        articleRepository.save(article);

        return null;
    }

    @Override
    public Article getArticleById(Long id) {
        return null;
    }

    @Override
    public Article updateArticle(Article article) {
        return null;
    }

    @Override
    public Boolean deleteArticleById(Long id) {
        return null;
    }
}
