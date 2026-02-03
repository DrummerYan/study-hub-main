package system

import (
	"context"

	sysModel "github.com/KeSilent/study-hub/server/model/system"
	"github.com/KeSilent/study-hub/server/service/system"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const initOrderMenuAuthority = initOrderMenu + initOrderAuthority

type initMenuAuthority struct{}

// auto run
func init() {
	system.RegisterInit(initOrderMenuAuthority, &initMenuAuthority{})
}

func (i *initMenuAuthority) MigrateTable(ctx context.Context) (context.Context, error) {
	return ctx, nil // do nothing
}

func (i *initMenuAuthority) TableCreated(ctx context.Context) bool {
	return false // always replace
}

func (i initMenuAuthority) InitializerName() string {
	return "sys_menu_authorities"
}

func (i *initMenuAuthority) InitializeData(ctx context.Context) (next context.Context, err error) {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return ctx, system.ErrMissingDBContext
	}
	authorities, ok := ctx.Value(initAuthority{}.InitializerName()).([]sysModel.SysAuthority)
	if !ok {
		return ctx, errors.Wrap(system.ErrMissingDependentContext, "创建 [菜单-权限] 关联失败, 未找到权限表初始化数据")
	}
	menus, ok := ctx.Value(initMenu{}.InitializerName()).([]sysModel.SysBaseMenu)
	if !ok {
		return next, errors.Wrap(errors.New(""), "创建 [菜单-权限] 关联失败, 未找到菜单表初始化数据")
	}
	next = ctx

	authorityByID := map[uint]*sysModel.SysAuthority{}
	for i := range authorities {
		authorityByID[authorities[i].AuthorityId] = &authorities[i]
	}

	pickMenus := func(paths ...string) []sysModel.SysBaseMenu {
		if len(paths) == 0 {
			return nil
		}
		menuByPath := map[string]sysModel.SysBaseMenu{}
		for _, menu := range menus {
			menuByPath[menu.Path] = menu
		}
		seen := map[string]struct{}{}
		picked := make([]sysModel.SysBaseMenu, 0, len(paths))
		for _, path := range paths {
			if _, exists := seen[path]; exists {
				continue
			}
			if menu, ok := menuByPath[path]; ok {
				picked = append(picked, menu)
				seen[path] = struct{}{}
			}
		}
		return picked
	}

	// 888 (超级管理员) - 全部菜单
	if auth := authorityByID[888]; auth != nil {
		if err = db.Model(auth).Association("SysBaseMenus").Replace(menus); err != nil {
			return next, err
		}
	}

	// 9001 (教师) - 课程/报名/课时/组织
	if auth := authorityByID[9001]; auth != nil {
		teacherMenus := pickMenus("dashboard", "person", "eduOrganization", "eduCourse", "eduEnrollment", "eduClassSession")
		if err = db.Model(auth).Association("SysBaseMenus").Replace(teacherMenus); err != nil {
			return next, err
		}
	}

	// 9002 (学员) - 只读入口
	if auth := authorityByID[9002]; auth != nil {
		studentMenus := pickMenus("dashboard", "person", "eduCourse", "eduClassSession")
		if err = db.Model(auth).Association("SysBaseMenus").Replace(studentMenus); err != nil {
			return next, err
		}
	}

	// 9003 (教务管理员/财务) - 课程/报名/课时/组织
	if auth := authorityByID[9003]; auth != nil {
		managerMenus := pickMenus("dashboard", "person", "eduOrganization", "eduCourse", "eduEnrollment", "eduClassSession")
		if err = db.Model(auth).Association("SysBaseMenus").Replace(managerMenus); err != nil {
			return next, err
		}
	}

	return next, nil
}

func (i *initMenuAuthority) DataInserted(ctx context.Context) bool {
	db, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return false
	}
	auth := &sysModel.SysAuthority{}
	if ret := db.Model(auth).
		Where("authority_id = ?", 888).Preload("SysBaseMenus").Find(auth); ret != nil {
		if ret.Error != nil {
			return false
		}
		return len(auth.SysBaseMenus) > 0
	}
	return false
}
