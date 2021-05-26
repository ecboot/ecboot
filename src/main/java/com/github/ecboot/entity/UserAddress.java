package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class UserAddress {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long addressId;

    private String addressName;

    private Long userId;

    private String consignee;

    private String email;

    private Long country;

    private Long province;

    private Long city;

    private Long district;

    private String address;

    private String zipcode;

    private String tel;

    private String mobile;

    private String signBuilding;

    private String bestTime;

}
