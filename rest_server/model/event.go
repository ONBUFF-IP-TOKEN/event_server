package model

import (
	"database/sql"
	"fmt"

	"github.com/ONBUFF-IP-TOKEN/baseutil/datetime"
	"github.com/ONBUFF-IP-TOKEN/baseutil/log"
	"github.com/ONBUFF-IP-TOKEN/event-server/rest_server/controllers/context"
)

func (o *DB) GetEventInfo(walletAddr string) (*context.Submit, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM event_attendees WHERE wallet_address='%v'", walletAddr)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	var ret sql.NullString
	var submitCnt sql.NullInt64
	info := &context.Submit{}
	if rows.Next() {
		if err := rows.Scan(&info.Idx, &info.WalletAddr, &info.ItemNum, &info.Email, &info.Ts, &ret, &submitCnt, &info.LastBalance); err != nil {
			log.Error(err)
		}
		info.Ret = ret.String
		info.SubmitCnt = submitCnt.Int64
	}
	return info, nil
}

func (o *DB) PostResetPurchase(resetPurchase *context.ResetPurchase) error {
	sqlQuery := "UPDATE event_item set owner=?, purchase_tx_hash=?, purchase_ts=? WHERE idx=?"

	_, err := o.Mysql.PrepareAndExec(sqlQuery, "", "", nil, resetPurchase.ItemNum)
	if err != nil {
		log.Error(err)
		return err
	}

	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (o *DB) PutEventSubmit(submit *context.Submit) (int64, error) {
	sqlQuery := fmt.Sprintf("INSERT INTO event_attendees(wallet_address, item_num, email, ts, submit_cnt, balance) VALUES (?,?,?,?,?,?)")

	result, err := o.Mysql.PrepareAndExec(sqlQuery, submit.WalletAddr, submit.ItemNum, submit.Email, submit.Ts, submit.SubmitCnt, submit.LastBalance)
	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("insert id:", insertId)
	return insertId, nil
}

func (o *DB) UpdateEventSubmit(submit *context.Submit) (int64, error) {
	sqlQuery := "UPDATE event_attendees set ts=?, submit_cnt=?, balance=? WHERE wallet_address=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, submit.Ts, submit.SubmitCnt, submit.LastBalance, submit.WalletAddr)
	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	log.Debug("update id:", insertId)
	return insertId, nil
}

func (o *DB) UpdateEventVerifyPurchase(purchaseInfo *context.PurchaseNoti) (int64, error) {
	sqlQuery := "UPDATE event_item set owner=?, purchase_tx_hash=?, purchase_ts=? WHERE idx=?"

	result, err := o.Mysql.PrepareAndExec(sqlQuery, purchaseInfo.WalletAddr, purchaseInfo.PurchaseTxHash, datetime.GetTS2MilliSec(), purchaseInfo.ItemNum)
	if err != nil {
		log.Error(err)
		return -1, err
	}
	insertId, err := result.LastInsertId()
	if err != nil {
		log.Error(err)
		return -1, err
	}
	return insertId, nil
}

func (o *DB) GetEventItem(itemNum int64) (*EventItemInfo, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM event_item WHERE idx=%v", itemNum)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	item := &EventItemInfo{}

	var owner, purchaseTxHash, tokenUri, info sql.NullString
	var tokenId, purchaseTs, submitStart, submitEnd, minAmount sql.NullInt64
	if rows.Next() {
		if err := rows.Scan(&item.Idx, &item.Name, &item.Serial, &tokenId, &tokenUri, &owner,
			&purchaseTxHash, &purchaseTs, &submitStart, &submitEnd, &minAmount, &item.Price, &info); err != nil {
			log.Error(err)
		}
		item.Owner = owner.String
		item.TokenId = tokenId.Int64
		item.TokenUri = tokenUri.String
		item.PurchaseTxHash = purchaseTxHash.String
		item.PurchaseTs = purchaseTs.Int64
		item.SubmitStart = submitStart.Int64
		item.SubmitEnd = submitEnd.Int64
		item.MinAmountForSumbit = minAmount.Int64
		item.Info = info.String
	}
	return item, nil
}

func (o *DB) GetLatestSubmitList(itemIdx int64) ([]context.Submit, error) {
	sqlQuery := fmt.Sprintf("SELECT * FROM event_attendees WHERE item_num=%v ORDER BY ts DESC LIMIT 0, 5", itemIdx)
	rows, err := o.Mysql.Query(sqlQuery)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer rows.Close()

	submits := make([]context.Submit, 0)
	var ret sql.NullString
	var submitCnt sql.NullInt64

	for rows.Next() {
		info := context.Submit{}
		if err := rows.Scan(&info.Idx, &info.WalletAddr, &info.ItemNum, &info.Email, &info.Ts, &ret, &submitCnt, &info.LastBalance); err != nil {
			log.Error(err)
		}
		info.Email = ""
		info.SubmitCnt = submitCnt.Int64

		first := info.WalletAddr[0:6]
		last := info.WalletAddr[len(info.WalletAddr)-4:]
		hide := ""
		for i := 0; i < len(info.WalletAddr)-6-4; i++ {
			hide = hide + "*"
		}
		info.WalletAddr = first + hide + last

		submits = append(submits, info)
	}
	return submits, nil
}
