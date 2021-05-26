package com.github.ecboot.entity;

import com.github.ecboot.enums.GenderEnum;
import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import java.util.Date;

@Entity
@Data
public class User {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    private String username;
    private String password;

    private String firstName;
    private String lastName;
    private String avatar;
    private GenderEnum gender;
    private Date birthday;
    private String motto;

    private String email;
    private Date emailVerifiedTime;
    private String mobile;
    private Date mobileVerifiedTime;

    private Integer locked;

    private Date CreatedTime;
    private Date UpdatedTime;
    private Date DeletedTime;
}
