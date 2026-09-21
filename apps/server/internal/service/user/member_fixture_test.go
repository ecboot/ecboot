package user

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

// ---- 会员中心（011）测试 fixture（自建+清理; init 基座见 auth_flow_test.go） ----

// seedMember 建测试会员（手机号密文+哈希, 含成长值）, 返回 userId。
func seedMember(ctx context.Context, t *gtest.T, phone, nickname string, growth int) int64 {
	mp := phoneCipher()
	hash := mp.Hash(phone)
	_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", hash)
	enc, err := mp.Encrypt(phone)
	t.AssertNil(err)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO `user`(nickname,phone,phone_hash,growth_value,status) VALUES(?,?,?,?,1)",
		nickname, enc, hash, growth)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// cleanupMember 清理会员及其登录日志/偏好。
func cleanupMember(ctx context.Context, t *gtest.T, phone string) {
	h := phoneCipher().Hash(phone)
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_login_log WHERE user_id IN (SELECT id FROM `user` WHERE phone_hash=?)", h)
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_notify_preference WHERE user_id IN (SELECT id FROM `user` WHERE phone_hash=?)", h)
	_, _ = g.DB().Exec(ctx, "DELETE FROM `user` WHERE phone_hash=?", h)
}

// seedLevelRule 建等级规则（用后 cleanupLevelRule）。
func seedLevelRule(ctx context.Context, t *gtest.T, name string, threshold int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_level_rule WHERE name=?", name)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO user_level_rule(name,growth_threshold,status) VALUES(?,?,1)", name, threshold)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupLevelRule(ctx context.Context, t *gtest.T, names ...string) {
	for _, n := range names {
		_, _ = g.DB().Exec(ctx, "DELETE FROM user_level_rule WHERE name=?", n)
	}
}

// seedLoginLog 建登录日志（used 于登录记录查询断言）。
func seedLoginLog(ctx context.Context, t *gtest.T, userId int64, daysAgo int, channel, status int) {
	_, err := g.DB().Exec(ctx,
		"INSERT INTO user_login_log(user_id,login_channel,login_status,ip,user_agent,created_at) "+
			"VALUES(?,?,?,'127.0.0.1','UA-test', DATE_SUB(NOW(), INTERVAL ? DAY))",
		userId, channel, status, daysAgo)
	t.AssertNil(err)
}

// seedSpuForMember 建测试 SPU（会员域自包含 fixture——跨包不能复用 shop 的测试 helper）。
func seedSpuForMember(ctx context.Context, t *gtest.T, suffix, price, image string, status int) int64 {
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE spu_no=?", "UM-SPU-"+suffix)
	res, err := g.DB().Exec(ctx,
		"INSERT INTO product_spu(spu_no,name,category_id,brand_id,images,price_min,price_max,status) VALUES(?,?,1,1,?,?,?,?)",
		"UM-SPU-"+suffix, "UM商品"+suffix, image, price, price, status)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

func cleanupSpuForMember(ctx context.Context, t *gtest.T, suffix string) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM product_spu WHERE spu_no=?", "UM-SPU-"+suffix)
}

// seedMessage 建站内信（isRead=true 时同时回填 read_time）。
func seedMessage(ctx context.Context, t *gtest.T, userId int64, title string, isRead bool) int64 {
	read, readTime := 0, "NULL"
	if isRead {
		read, readTime = 1, "NOW()"
	}
	res, err := g.DB().Exec(ctx,
		"INSERT INTO user_message(user_id,title,content,biz_type,biz_no,is_read,read_time) "+
			"VALUES(?,?,'内容',1,'',?,"+readTime+")", userId, title, read)
	t.AssertNil(err)
	id, _ := res.LastInsertId()
	return id
}

// seedPointAccount 建积分账户（balance 可负; lastEarnedAt 为天数偏移, nil=不设）。
func seedPointAccount(ctx context.Context, t *gtest.T, userId int64, balance int, lastEarnedDaysAgo *int) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM point_account WHERE user_id=?", userId)
	le := "NULL"
	if lastEarnedDaysAgo != nil {
		le = "DATE_SUB(NOW(), INTERVAL " + strconv.Itoa(*lastEarnedDaysAgo) + " DAY)"
	}
	_, err := g.DB().Exec(ctx,
		"INSERT INTO point_account(user_id,balance,last_earned_at) VALUES(?,?,"+le+")", userId, balance)
	t.AssertNil(err)
}

func cleanupPoint(ctx context.Context, t *gtest.T, userId int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM point_log WHERE user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM point_account WHERE user_id=?", userId)
}

// seedPointLog 建积分流水。
func seedPointLog(ctx context.Context, t *gtest.T, userId int64, bizType, points, balanceAfter int, orderNo string) {
	_, err := g.DB().Exec(ctx,
		"INSERT INTO point_log(user_id,biz_type,points,balance_after,order_no) VALUES(?,?,?,?,?)",
		userId, bizType, points, balanceAfter, orderNo)
	t.AssertNil(err)
}

// seedInviteRecord 建邀请记录（inviter 邀请 newUser; 第三参为 **reward_trigger** 时机: 1注册即发 2首单后发）。
func seedInviteRecord(ctx context.Context, t *gtest.T, inviterId, newUserId int64, rewardType int) {
	// 按唯一键（new_user_id）清理——按 inviter 清会误删同一邀请人的其他记录
	_, _ = g.DB().Exec(ctx, "DELETE FROM invite_record WHERE new_user_id=?", newUserId)
	_, err := g.DB().Exec(ctx,
		"INSERT INTO invite_record(new_user_id,inviter_id,reward_type,reward_ref,status,reward_trigger) VALUES(?,?,1,0,1,?)",
		newUserId, inviterId, rewardType)
	t.AssertNil(err)
}

func cleanupInvite(ctx context.Context, t *gtest.T, inviterId int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM invite_record WHERE inviter_id=?", inviterId)
}

// cleanupMemberCollections 清理会员的收藏与足迹（评审 I4: 此前遗漏致孤儿行累积）。
func cleanupMemberCollections(ctx context.Context, t *gtest.T, userId int64) {
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_favorite WHERE user_id=?", userId)
	_, _ = g.DB().Exec(ctx, "DELETE FROM user_footprint WHERE user_id=?", userId)
}
