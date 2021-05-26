package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class Cart {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long recId;

    private Long userId;

    private String sessionId;

    private Long goodsId;

    private String goodsSn;

    private Long productId;

    private String goodsName;

    private BigDecimal marketPrice;

    private BigDecimal goodsPrice;

    private Long goodsNumber;

    private String goodsAttr;

    private Long isReal;

    private String extensionCode;

    private Long parentId;

    private Long recType;

    private Long isGift;

    private Long isShipping;

    private Long canHandsel;

    private String goodsAttrId;

}
