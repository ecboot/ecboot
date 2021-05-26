package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class UserBonus {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long bonusId;

    private Long bonusTypeId;

    private Long bonusSn;

    private Long userId;

    private Long usedTime;

    private Long orderId;

    private Long emailed;

}
