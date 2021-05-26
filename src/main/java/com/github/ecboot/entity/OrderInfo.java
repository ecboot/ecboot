package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class OrderInfo {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long orderId;

    private String orderSn;

    private Long userId;

    private Long orderStatus;

    private Long shippingStatus;

    private Long payStatus;

    private String consignee;

    private Long country;

    private Long province;

    private Long city;

    private Long district;

    private String address;

    private String zipcode;

    private String tel;

    private String mobile;

    private String email;

    private String bestTime;

    private String signBuilding;

    private String postscript;

    private Long shippingId;

    private String shippingName;

    private Long payId;

    private String payName;

    private String howOos;

    private String howSurplus;

    private String packName;

    private String cardName;

    private String cardMessage;

    private String invPayee;

    private String invContent;

    private BigDecimal goodsAmount;

    private BigDecimal shippingFee;

    private BigDecimal insureFee;

    private BigDecimal payFee;

    private BigDecimal packFee;

    private BigDecimal cardFee;

    private BigDecimal moneyPaid;

    private BigDecimal surplus;

    private Long integral;

    private BigDecimal integralMoney;

    private BigDecimal bonus;

    private BigDecimal orderAmount;

    private Long fromAd;

    private String referer;

    private Long addTime;

    private Long confirmTime;

    private Long payTime;

    private Long shippingTime;

    private Long packId;

    private Long cardId;

    private Long bonusId;

    private String invoiceNo;

    private String extensionCode;

    private Long extensionId;

    private String toBuyer;

    private String payNote;

    private Long agencyId;

    private String invType;

    private BigDecimal tax;

    private Long isSeparate;

    private Long parentId;

    private BigDecimal discount;

}
