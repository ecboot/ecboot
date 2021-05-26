package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class OrderGoods {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long recId;

    private Long orderId;

    private Long goodsId;

    private String goodsName;

    private String goodsSn;

    private Long productId;

    private Long goodsNumber;

    private BigDecimal marketPrice;

    private BigDecimal goodsPrice;

    private String goodsAttr;

    private Long sendNumber;

    private Long isReal;

    private String extensionCode;

    private Long parentId;

    private Long isGift;

    private String goodsAttrId;

}
