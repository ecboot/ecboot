package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Attribute {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long attrId;

    private Long catId;

    private String attrName;

    private Long attrInputType;

    private Long attrType;

    private String attrValues;

    private Long attrIndex;

    private Long sortOrder;

    private Long isLinked;

    private Long attrGroup;

}
