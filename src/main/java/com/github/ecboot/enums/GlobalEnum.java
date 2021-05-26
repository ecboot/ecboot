package com.github.ecboot.enums;

import lombok.Getter;

@Getter
public enum GlobalEnum {

    /* 图片处理相关常数 */
    ERR_INVALID_IMAGE(1, ""),
    ERR_NO_GD(2, ""),
    ERR_IMAGE_NOT_EXISTS(3, ""),
    ERR_DIRECTORY_READONLY(4, ""),
    ERR_UPLOAD_FAILURE(5, ""),
    ERR_INVALID_PARAM(6, ""),
    ERR_INVALID_IMAGE_TYPE(7, ""),

    /* 插件相关常数 */
    ERR_COPYFILE_FAILED(1, ""),
    ERR_CREATE_TABLE_FAILED(2, ""),
    ERR_DELETE_FILE_FAILED(3, ""),

    /* 商品属性类型常数 */
    ATTR_TEXT(0, ""),
    ATTR_OPTIONAL(1, ""),
    ATTR_TEXTAREA(2, ""),
    ATTR_URL(3, ""),

    /* 会员整合相关常数 */
    ERR_USERNAME_EXISTS(1, "用户名已经存在"),
    ERR_EMAIL_EXISTS(2, "Email已经存在"),
    ERR_INVALID_USER_ID(3, "无效的user_id"),
    ERR_INVALID_USERNAME(4, "无效的用户名"),
    ERR_INVALID_PASSWORD(5, "密码错误"),
    ERR_INVALID_EMAIL(6, "email错误"),
    ERR_USERNAME_NOT_ALLOW(7, "用户名不允许注册"),
    ERR_EMAIL_NOT_ALLOW(8, "EMAIL不允许注册"),

    /* 加入购物车失败的错误代码 */
    ERR_NOT_EXISTS(1, "商品不存在"),
    ERR_OUT_OF_STOCK(2, "商品缺货"),
    ERR_NOT_ON_SALE(3, "商品已下架"),
    ERR_CANNOT_ALONE_SALE(4, "商品不能单独销售"),
    ERR_NO_BASIC_GOODS(5, "没有基本件"),
    ERR_NEED_SELECT_ATTR(6, "需要用户选择属性"),

    /* 购物车商品类型 */
    CART_GENERAL_GOODS(0, "普通商品"),
    CART_GROUP_BUY_GOODS(1, "团购商品"),
    CART_AUCTION_GOODS(2, "拍卖商品"),
    CART_SNATCH_GOODS(3, "夺宝奇兵"),
    CART_EXCHANGE_GOODS(4, "积分商城"),

    /* 订单状态 */
    OS_UNCONFIRMED(0, "未确认"),
    OS_CONFIRMED(1, "已确认"),
    OS_CANCELED(2, "已取消"),
    OS_INVALID(3, "无效"),
    OS_RETURNED(4, "退货"),
    OS_SPLITED(5, "已分单"),
    OS_SPLITING_PART(6, "部分分单"),

    /* 支付类型 */
    PAY_ORDER(0, "订单支付"),
    PAY_SURPLUS(1, "会员预付款"),

    /* 配送状态 */
    SS_UNSHIPPED(0, "未发货"),
    SS_SHIPPED(1, "已发货"),
    SS_RECEIVED(2, "已收货"),
    SS_PREPARING(3, "备货中"),
    SS_SHIPPED_PART(4, "已发货(部分商品)"),
    SS_SHIPPED_ING(5, "发货中(处理分单)"),
    OS_SHIPPED_PART(6, "已发货(部分商品)"),

    /* 支付状态 */
    PS_UNPAYED(0, "未付款"),
    PS_PAYING(1, "付款中"),
    PS_PAYED(2, "已付款"),

    /* 综合状态 */
    CS_AWAIT_PAY(100, "待付款：货到付款且已发货且未付款，非货到付款且未付款"),
    CS_AWAIT_SHIP(101, "待发货：货到付款且未发货，非货到付款且已付款且未发货"),
    CS_FINISHED(102, "已完成：已确认、已付款、已发货"),

    /* 缺货处理 */
    OOS_WAIT(0, "等待货物备齐后再发"),
    OOS_CANCEL(1, "取消订单"),
    OOS_CONSULT(2, "与店主协商"),

    /* 帐户明细类型 */
    SURPLUS_SAVE(0, "为帐户冲值"),
    SURPLUS_RETURN(1, "从帐户提款"),

    /* 评论状态 */
    COMMENT_UNCHECKED(0, "未审核"),
    COMMENT_CHECKED(1, "已审核或已回复(允许显示)"),
    COMMENT_REPLYED(2, "该评论的内容属于回复"),

    /* 红包发放的方式 */
    SEND_BY_USER(0, "按用户发放"),
    SEND_BY_GOODS(1, "按商品发放"),
    SEND_BY_ORDER(2, "按订单发放"),
    SEND_BY_PRINT(3, "线下发放"),

    /* 广告的类型 */
    IMG_AD(0, "图片广告"),
    FALSH_AD(1, "flash广告"),
    CODE_AD(2, "代码广告"),
    TEXT_AD(3, "文字广告"),

    /* 是否需要用户选择属性 */
    ATTR_NOT_NEED_SELECT(0, "不需要选择"),
    ATTR_NEED_SELECT(1, "需要选择"),

    /* 用户中心留言类型 */
    M_MESSAGE(0, "留言"),
    M_COMPLAINT(1, "投诉"),
    M_ENQUIRY(2, "询问"),
    M_CUSTOME(3, "售后"),
    M_BUY(4, "求购"),
    M_BUSINESS(5, "商家"),
    M_COMMENT(6, "评论"),

    /* 团购活动状态 */
    GBS_PRE_START(0, "未开始"),
    GBS_UNDER_WAY(1, "进行中"),
    GBS_FINISHED(2, "已结束"),
    GBS_SUCCEED(3, "团购成功（可以发货了）"),
    GBS_FAIL(4, "团购失败"),

    /* 红包是否发送邮件 */
    BONUS_NOT_MAIL(0, "不发送"),
    BONUS_MAIL_SUCCEED(1, "发送成功"),
    BONUS_MAIL_FAIL(2, "发送失败"),

    /* 商品活动类型 */
    GAT_SNATCH(0, ""),
    GAT_GROUP_BUY(1, ""),
    GAT_AUCTION(2, ""),
    GAT_POINT_BUY(3, ""),
    GAT_PACKAGE(4, "超值礼包"),

    /* 帐号变动类型 */
    ACT_SAVING(0, "帐户冲值"),
    ACT_DRAWING(1, "帐户提款"),
    ACT_ADJUSTING(2, "调节帐户"),
    ACT_OTHER(99, "其他类型"),

    /* 密码加密方法 */
    PWD_MD5(1, "md5加密方式"),
    PWD_PRE_SALT(2, "前置验证串的加密方式"),
    PWD_SUF_SALT(3, "后置验证串的加密方式"),

    /* 文章分类类型 */
    COMMON_CAT(1, "普通分类"),
    SYSTEM_CAT(2, "系统默认分类"),
    INFO_CAT(3, "网店信息分类"),
    UPHELP_CAT(4, "网店帮助分类分类"),
    HELP_CAT(5, "网店帮助分类"),

    /* 活动状态 */
    PRE_START(0, "未开始"),
    UNDER_WAY(1, "进行中"),
    FINISHED(2, "已结束"),
    SETTLED(3, "已处理"),

    /* 验证码 */
    CAPTCHA_REGISTER(1, "注册时使用验证码"),
    CAPTCHA_LOGIN(2, "登录时使用验证码"),
    CAPTCHA_COMMENT(4, "评论时使用验证码"),
    CAPTCHA_ADMIN(8, "后台登录时使用验证码"),
    CAPTCHA_LOGIN_FAIL(16, "登录失败后显示验证码"),
    CAPTCHA_MESSAGE(32, "留言时使用验证码"),

    /* 优惠活动的优惠范围 */
    FAR_ALL(0, "全部商品"),
    FAR_CATEGORY(1, "按分类选择"),
    FAR_BRAND(2, "按品牌选择"),
    FAR_GOODS(3, "按商品选择"),

    /* 优惠活动的优惠方式 */
    FAT_GOODS(0, "送赠品或优惠购买"),
    FAT_PRICE(1, "现金减免"),
    FAT_DISCOUNT(2, "价格打折优惠"),

    /* 评论条件 */
    COMMENT_LOGIN(1, "只有登录用户可以评论"),
    COMMENT_CUSTOM(2, "只有有过一次以上购买行为的用户可以评论"),
    COMMENT_BOUGHT(3, "只有购买够该商品的人可以评论"),

    /* 减库存时机 */
    SDT_SHIP(0, "发货时"),
    SDT_PLACE(1, "下订单时"),

    /* 加密方式 */
    ENCRYPT_ZC(1, "zc加密方式"),
    ENCRYPT_UC(2, "uc加密方式"),

    /* 商品类别 */
    G_REAL(1, "实体商品"),
    G_CARD(0, "虚拟卡"),

    /* 积分兑换 */
    TO_P(0, "兑换到商城消费积分"),
    FROM_P(1, "用商城消费积分兑换"),
    TO_R(2, "兑换到商城等级积分"),
    FROM_R(3, "用商城等级积分兑换"),

    /* 添加feed事件到UC的TYPE*/
    BUY_GOODS(1, "购买商品"),
    COMMENT_GOODS(2, "添加商品评论"),

    /* 邮件发送用户 */
    SEND_LIST(0, ""),
    SEND_USER(1, ""),
    SEND_RANK(2, ""),
    ;

    private final Integer code;

    private final String description;

    GlobalEnum(Integer code, String description) {
        this.code = code;
        this.description = description;
    }
}
