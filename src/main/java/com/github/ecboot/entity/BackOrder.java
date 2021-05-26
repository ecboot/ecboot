package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class BackOrder {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long backId;

    private String deliverySn;

    private String orderSn;

    private Long orderId;

    private String invoiceNo;

    private Long addTime;

    private Long shippingId;

    private String shippingName;

    private Long userId;

    private String actionUser;

    private String consignee;

    private String address;

    private Long country;

    private Long province;

    private Long city;

    private Long district;

    private String signBuilding;

    private String email;

    private String zipcode;

    private String tel;

    private String mobile;

    private String bestTime;

    private String postscript;

    private String howOos;

    private BigDecimal insureFee;

    private BigDecimal shippingFee;

    private Long updateTime;

    private Long suppliersId;

    private Long status;

    private Long returnTime;

    private Long agencyId;

}
