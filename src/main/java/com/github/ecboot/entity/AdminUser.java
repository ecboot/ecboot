package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class AdminUser {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long userId;

    private String userName;

    private String email;

    private String password;

    private String ecSalt;

    private Long addTime;

    private Long lastLogin;

    private String lastIp;

    private String actionList;

    private String navList;

    private String langType;

    private Long agencyId;

    private Long suppliersId;

    private String todolist;

    private Long roleId;

}
