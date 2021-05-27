package com.github.ecboot.manager.sender;

import com.github.ecboot.manager.sender.sms.AliyunAdaptor;

public class SmsSender implements Sender {

    public void send() {
        AliyunAdaptor aliyunAdaptor = new AliyunAdaptor();
        aliyunAdaptor.send();
    }
}  