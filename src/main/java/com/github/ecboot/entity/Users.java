package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.math.BigDecimal;
import java.util.Date;

@Entity
@Data
public class Users {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long userId;

    private String email;

    private String userName;

    private String password;

    private String question;

    private String answer;

    private Long sex;

    private Date birthday;

    private BigDecimal userMoney;

    private BigDecimal frozenMoney;

    private Long payPoints;

    private Long rankPoints;

    private Long addressId;

    private Long regTime;

    private Long lastLogin;

    private Date lastTime;

    private String lastIp;

    private Long visitCount;

    private Long userRank;

    private Long isSpecial;

    private String ecSalt;

    private String salt;

    private Long parentId;

    private Long flag;

    private String alias;

    private String msn;

    private String qq;

    private String officePhone;

    private String homePhone;

    private String mobilePhone;

    private Long isValidated;

    private BigDecimal creditLine;

    private String passwdQuestion;

    private String passwdAnswer;

}
