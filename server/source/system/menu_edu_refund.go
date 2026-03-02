package system

import (
	"context"
	"errors"

	sysModel "github.com/KeSilent/study-hub/server/model/system"
	"github.com/KeSilent/study-hub/server/service/system"
	"gorm.io/gorm"
)

const initOrderEduRefundMenu = initOrderMenuAuthority + 1

type initEduRefundMenu struct{}

// auto run
func init() {
	system.RegisterInit(initOrderEduRefundMenu, &initEduRefundMenu{})
}

func (i initEduRefundMenu) InitializerName() string {
	return "edu_refund_menu"
}

func (i *initEduRefundMenu) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func (i *initEduRefundMenu) TableCreated(ctx context.Context) bool {
	return true
}

func (i *initEduRefundMenu) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	if errors.Is(db.Where("path = ?", "eduRefund").First(&sysModel.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (i *initEduRefundMenu) InitializeData(ctx context.Context) (context.Context, error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	var existing sysModel.SysBaseMenu
	if err := db.Where("path = ?", "eduRefund").First(&existing).Error; err == nil {
		return ctx, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return ctx, err
	}

	menu := sysModel.SysBaseMenu{
		MenuLevel: 0,
		Hidden:    false,
		ParentId:  "0",
		Path:      "eduRefund",
		Name:      "eduRefund",
		Component: "view/eduRefund/eduRefund.vue",
		Sort:      5,
		Meta:      sysModel.Meta{Title: "退费记录", Icon: "money"},
	}
	if err := db.Create(&menu).Error; err != nil {
		return ctx, err
	}

	var auth sysModel.SysAuthority
	if err := db.Where("authority_id = ?", 888).First(&auth).Error; err == nil {
		_ = db.Model(&auth).Association("SysBaseMenus").Append(&menu)
	}

	return ctx, nil
}
