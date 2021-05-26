package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Article {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long articleId;

    private Long catId;

    private String title;

    private String content;

    private String author;

    private String authorEmail;

    private String keywords;

    private Long articleType;

    private Long isOpen;

    private Long addTime;

    private String fileUrl;

    private Long openType;

    private String link;

    private String description;

}
