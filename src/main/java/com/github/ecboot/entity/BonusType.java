package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class BonusType {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long typeId;

    private String typeName;

    private BigDecimal typeMoney;

    private Long sendType;

    private BigDecimal minAmount;

    private BigDecimal maxAmount;

    private Long sendStartDate;

    private Long sendEndDate;

    private Long useStartDate;

    private Long useEndDate;

    private BigDecimal minGoodsAmount;

}
