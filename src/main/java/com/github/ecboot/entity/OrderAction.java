package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class OrderAction {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long actionId;

    private Long orderId;

    private String actionUser;

    private Long orderStatus;

    private Long shippingStatus;

    private Long payStatus;

    private Long actionPlace;

    private String actionNote;

    private Long logTime;

}
