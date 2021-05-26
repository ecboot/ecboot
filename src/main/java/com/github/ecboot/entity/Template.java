package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Template {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id; // TODO

    private String filename;

    private String region;

    private String library;

    private Long sortOrder;

    private Long number;

    private Long type;

    private String theme;

    private String remarks;

}
