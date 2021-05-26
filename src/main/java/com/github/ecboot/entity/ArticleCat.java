package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class ArticleCat {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long catId;

    private String catName;

    private Long catType;

    private String keywords;

    private String catDesc;

    private Long sortOrder;

    private Long showInNav;

    private Long parentId;

}
