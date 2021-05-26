package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class DeliveryGoods {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long recId;

    private Long deliveryId;

    private Long goodsId;

    private Long productId;

    private String productSn;

    private String goodsName;

    private String brandName;

    private String goodsSn;

    private Long isReal;

    private String extensionCode;

    private Long parentId;

    private Long sendNumber;

    private String goodsAttr;

}
