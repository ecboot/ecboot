package com.github.ecboot.service.article;

import com.github.ecboot.entity.Article;

import java.util.List;

public interface ArticleService {

    List<Article> getList();

    Article create(Article article);

    Article getArticleById(Long id);

    Article updateArticle(Article article);

    Boolean deleteArticleById(Long id);
}
