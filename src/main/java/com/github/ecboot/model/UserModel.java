package com.github.ecboot.model;

import com.github.ecboot.enums.GenderEnum;
import lombok.Data;

import java.util.Date;

@Data
public class UserModel {

    private String firstName;
    private String lastName;
    private String avatar;
    private GenderEnum gender;
    private Date birthday;
    private String motto;
}
