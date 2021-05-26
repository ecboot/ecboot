package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Shipping {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long shippingId;

    private String shippingCode;

    private String shippingName;

    private String shippingDesc;

    private String insure;

    private Long supportCod;

    private Long enabled;

    private String shippingPrint;

    private String printBg;

    private String configLable;

    private Long printModel;

    private Long shippingOrder;

}
