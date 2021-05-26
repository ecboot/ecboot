package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Ad {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long adId;

    private Long positionId;

    private Long mediaType;

    private String adName;

    private String adLink;

    private String adCode;

    private Long startTime;

    private Long endTime;

    private String linkMan;

    private String linkEmail;

    private String linkPhone;

    private Long clickCount;

    private Long enabled;

}
