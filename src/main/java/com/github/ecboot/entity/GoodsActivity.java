package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class GoodsActivity {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long actId;

    private String actName;

    private String actDesc;

    private Long actType;

    private Long goodsId;

    private Long productId;

    private String goodsName;

    private Long startTime;

    private Long endTime;

    private Long isFinished;

    private String extInfo;

}
