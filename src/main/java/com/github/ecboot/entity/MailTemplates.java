package com.github.ecboot.entity;

import lombok.Data;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
@Data
public class MailTemplates {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long templateId;

    private String templateCode;

    private Long isHtml;

    private String templateSubject;

    private String templateContent;

    private Long lastModify;

    private Long lastSend;

    private String type;

}
