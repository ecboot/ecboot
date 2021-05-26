package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Stats {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id; // TODO

    private Long accessTime;

    private String ipAddress;

    private Long visitTimes;

    private String browser;

    private String system;

    private String language;

    private String area;

    private String refererDomain;

    private String refererPath;

    private String accessUrl;

}
