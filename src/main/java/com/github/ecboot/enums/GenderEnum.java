package com.github.ecboot.enums;

import lombok.Getter;

@Getter
public enum GenderEnum {

    UNKNOWN(0, "保密"),
    FEMALE(1, "女性"),
    MALE(2, "男性"),
    ;

    private final Integer code;

    private final String note;

    GenderEnum(Integer code, String note) {
        this.code = code;
        this.note = note;
    }

}
