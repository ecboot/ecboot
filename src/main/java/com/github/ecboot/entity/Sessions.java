package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class Sessions {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private String sesskey;

    private Long expiry;

    private Long userid;

    private Long adminid;

    private String ip;

    private String userName;

    private Long userRank;

    private BigDecimal discount;

    private String email;

    private String data;

}
