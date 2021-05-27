package com.github.ecboot.manager.sender;

public class SendMailFactory implements Provider {

    public Sender produce() {
        return new MailSender();
    }
}  