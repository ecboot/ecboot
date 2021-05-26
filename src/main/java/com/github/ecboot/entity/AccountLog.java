package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class AccountLog {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long logId;

    private Long userId;

    private BigDecimal userMoney;

    private BigDecimal frozenMoney;

    private Long rankPoints;

    private Long payPoints;

    private Long changeTime;

    private String changeDesc;

    private Long changeType;

}
