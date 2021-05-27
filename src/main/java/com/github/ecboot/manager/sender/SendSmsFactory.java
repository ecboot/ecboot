package com.github.ecboot.manager.sender;

public class SendSmsFactory implements Provider {

    public Sender produce() {
        return new SmsSender();
    }
}