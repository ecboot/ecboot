package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class BookingGoods {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long recId;

    private Long userId;

    private String email;

    private String linkMan;

    private String tel;

    private Long goodsId;

    private String goodsDesc;

    private Long goodsNumber;

    private Long bookingTime;

    private Long isDispose;

    private String disposeUser;

    private Long disposeTime;

    private String disposeNote;

}
