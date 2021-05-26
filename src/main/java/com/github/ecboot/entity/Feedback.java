package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Feedback {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long msgId;

    private Long parentId;

    private Long userId;

    private String userName;

    private String userEmail;

    private String msgTitle;

    private Long msgType;

    private Long msgStatus;

    private String msgContent;

    private Long msgTime;

    private String messageImg;

    private Long orderId;

    private Long msgArea;

}
