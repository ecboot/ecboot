package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Plugins {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private String code;

    private String version;

    private String library;

    private Long assign;

    private Long installDate;

}
