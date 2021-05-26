package com.github.ecboot.response;

import com.github.ecboot.enums.GenderEnum;
import lombok.Data;

import java.util.Date;

@Data
public class UserResponse {

    private String firstName;
    private String lastName;
    private String avatar;
    private GenderEnum gender;
    private Date birthday;
    private String motto;
}
