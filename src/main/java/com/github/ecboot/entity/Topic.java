package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Topic {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long topicId;

    private String title;

    private String intro;

    private Long startTime;

    private Long endTime;

    private String data;

    private String template;

    private String css;

    private String topicImg;

    private String titlePic;

    private String baseStyle;

    private String htmls;

    private String keywords;

    private String description;

}
