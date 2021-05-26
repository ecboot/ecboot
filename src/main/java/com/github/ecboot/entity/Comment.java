package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class Comment {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long commentId;

    private Long commentType;

    private Long idValue;

    private String email;

    private String userName;

    private String content;

    private Long commentRank;

    private Long addTime;

    private String ipAddress;

    private Long status;

    private Long parentId;

    private Long userId;

}
