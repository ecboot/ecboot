package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class VirtualCard {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long cardId;

    private Long goodsId;

    private String cardSn;

    private String cardPassword;

    private Long addDate;

    private Long endDate;

    private Long isSaled;

    private String orderSn;

    private String crc32;

}
