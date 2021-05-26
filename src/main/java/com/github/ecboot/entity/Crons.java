package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Crons {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long cronId;

    private String cronCode;

    private String cronName;

    private String cronDesc;

    private Long cronOrder;

    private String cronConfig;

    private Long thistime;

    private Long nextime;

    private Long day;

    private String week;

    private String hour;

    private String minute;

    private Long enable;

    private Long runOnce;

    private String allowIp;

    private String alowFiles;

}
