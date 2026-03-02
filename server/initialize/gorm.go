package initialize

import (
	"errors"
	"os"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/example"
	"github.com/KeSilent/study-hub/server/model/system"

	"github.com/KeSilent/study-hub/server/model/edu_organization"
	"github.com/KeSilent/study-hub/server/model/edu_user_course"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func Gorm() *gorm.DB {
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	case "pgsql":
		return GormPgSql()
	case "oracle":
		return GormOracle()
	case "mssql":
		return GormMssql()
	default:
		return GormMysql()
	}
}

func RegisterTables() {
	db := global.GVA_DB
	err := db.AutoMigrate(

		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.JwtBlacklist{},
		system.SysAuthority{},
		system.SysDictionary{},
		system.SysOperationRecord{},
		system.SysAutoCodeHistory{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},
		system.SysAutoCode{},
		system.SysChatGptOption{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},

		edu_user_course.EduClassSession{},
		edu_organization.EduCourse{},
		edu_organization.EduOrganization{},
		edu_user_course.EduEnrollment{},
		edu_user_course.EduPayment{},
		edu_user_course.EduRefund{},
	)
	if err != nil {
		global.GVA_LOG.Error("register table failed", zap.Error(err))
		os.Exit(0)
	}
	ensureEduRefundMenu(db)
	global.GVA_LOG.Info("register table success")
}

func ensureEduRefundMenu(db *gorm.DB) {
	if db == nil {
		return
	}
	var menu system.SysBaseMenu
	err := db.Where("path = ?", "eduRefund").First(&menu).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		menu = system.SysBaseMenu{
			MenuLevel: 0,
			Hidden:    false,
			ParentId:  "0",
			Path:      "eduRefund",
			Name:      "eduRefund",
			Component: "view/eduRefund/eduRefund.vue",
			Sort:      5,
			Meta:      system.Meta{Title: "退费记录", Icon: "money"},
		}
		if err := db.Create(&menu).Error; err != nil {
			global.GVA_LOG.Error("create refund menu failed", zap.Error(err))
			return
		}
	} else if err != nil {
		global.GVA_LOG.Error("query refund menu failed", zap.Error(err))
		return
	}

	var auth system.SysAuthority
	if err := db.Where("authority_id = ?", 888).First(&auth).Error; err != nil {
		return
	}
	_ = db.Model(&auth).Association("SysBaseMenus").Append(&menu)
}
