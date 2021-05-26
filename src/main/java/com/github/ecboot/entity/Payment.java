package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Payment {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long payId;

    private String payCode;

    private String payName;

    private String payFee;

    private String payDesc;

    private Long payOrder;

    private String payConfig;

    private Long enabled;

    private Long isCod;

    private Long isOnline;

}
