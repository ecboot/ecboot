package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;

@Entity
@Data
public class UserAccount {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private Long userId;

    private String adminUser;

    private BigDecimal amount;

    private Long addTime;

    private Long paidTime;

    private String adminNote;

    private String userNote;

    private Long processType;

    private String payment;

    private Long isPaid;

}
