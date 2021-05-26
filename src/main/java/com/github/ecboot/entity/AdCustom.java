package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class AdCustom {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long adId;

    private Long adType;

    private String adName;

    private Long addTime;

    private String content;

    private String url;

    private Long adStatus;

}
