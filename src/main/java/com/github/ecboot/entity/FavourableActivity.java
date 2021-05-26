package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class FavourableActivity {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long actId;

    private String actName;

    private Long startTime;

    private Long endTime;

    private String userRank;

    private Long actRange;

    private String actRangeExt;

    private BigDecimal minAmount;

    private BigDecimal maxAmount;

    private Long actType;

    private BigDecimal actTypeExt;

    private String gift;

    private Long sortOrder;

}
