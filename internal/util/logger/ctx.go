package logger

import "context"

type logCtx struct {
	UserID      int
	Username    string
	ToUser      string
	SendAmount  int
	CoinBalance int
	Item        string
}

type keyType int

const key = keyType(0)

func WithLogUserID(ctx context.Context, userID int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.UserID = userID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{UserID: userID})
}

func WithLogUsername(ctx context.Context, username string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Username = username
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Username: username})
}

func WithLogToUser(ctx context.Context, toUser string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.ToUser = toUser
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{ToUser: toUser})
}

func WithLogSendAmount(ctx context.Context, sendAmount int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.SendAmount = sendAmount
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{SendAmount: sendAmount})
}

func WithLogCoinBalance(ctx context.Context, coinBalance int) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.CoinBalance = coinBalance
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{CoinBalance: coinBalance})
}
func WithLogItem(ctx context.Context, item string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Item = item
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Item: item})
}
