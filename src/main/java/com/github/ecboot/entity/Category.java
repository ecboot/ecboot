package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Category {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long catId;

    private String catName;

    private String keywords;

    private String catDesc;

    private Long parentId;

    private Long sortOrder;

    private String templateFile;

    private String measureUnit;

    private Long showInNav;

    private String style;

    private Long isShow;

    private Long grade;

    private String filterAttr;

}
