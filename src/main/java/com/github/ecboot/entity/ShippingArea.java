package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class ShippingArea {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long shippingAreaId;

    private String shippingAreaName;

    private Long shippingId;

    private String configure;

}
