package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class AffiliateLog {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long logId;

    private Long orderId;

    private Long time;

    private Long userId;

    private String userName;

    private BigDecimal money;

    private Long point;

    private Long separateType;

}
