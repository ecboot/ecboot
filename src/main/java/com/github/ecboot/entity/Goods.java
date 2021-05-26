package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class Goods {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long goodsId;

    private Long catId;

    private String goodsSn;

    private String goodsName;

    private String goodsNameStyle;

    private Long clickCount;

    private Long brandId;

    private String providerName;

    private Long goodsNumber;

    private BigDecimal goodsWeight;

    private BigDecimal marketPrice;

    private BigDecimal shopPrice;

    private BigDecimal promotePrice;

    private Long promoteStartDate;

    private Long promoteEndDate;

    private Long warnNumber;

    private String keywords;

    private String goodsBrief;

    private String goodsDesc;

    private String goodsThumb;

    private String goodsImg;

    private String originalImg;

    private Long isReal;

    private String extensionCode;

    private Long isOnSale;

    private Long isAloneSale;

    private Long isShipping;

    private Long integral;

    private Long addTime;

    private Long sortOrder;

    private Long isDelete;

    private Long isBest;

    private Long isNew;

    private Long isHot;

    private Long isPromote;

    private Long bonusTypeId;

    private Long lastUpdate;

    private Long goodsType;

    private String sellerNote;

    private Long giveIntegral;

    private Long rankIntegral;

    private Long suppliersId;

    private Long isCheck;

}
