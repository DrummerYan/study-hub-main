package system

import (
	"errors"
	"strconv"

	"github.com/KeSilent/study-hub/server/global"
	"github.com/KeSilent/study-hub/server/model/common/request"
	"github.com/KeSilent/study-hub/server/model/system"
	"github.com/KeSilent/study-hub/server/model/system/response"
	"gorm.io/gorm"
)

var ErrRoleExistence = errors.New("存在相同角色id")

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CreateAuthority
//@description: 创建一个角色
//@param: auth model.SysAuthority
//@return: authority system.SysAuthority, err error

type AuthorityService struct{}

var AuthorityServiceApp = new(AuthorityService)

func (authorityService *AuthorityService) CreateAuthority(auth system.SysAuthority) (authority system.SysAuthority, err error) {
	var authorityBox system.SysAuthority
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return auth, ErrRoleExistence
	}
	err = global.GVA_DB.Create(&auth).Error
	return auth, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: CopyAuthority
//@description: 复制一个角色
//@param: copyInfo response.SysAuthorityCopyResponse
//@return: authority system.SysAuthority, err error

func (authorityService *AuthorityService) CopyAuthority(copyInfo response.SysAuthorityCopyResponse) (authority system.SysAuthority, err error) {
	var authorityBox system.SysAuthority
	if !errors.Is(global.GVA_DB.Where("authority_id = ?", copyInfo.Authority.AuthorityId).First(&authorityBox).Error, gorm.ErrRecordNotFound) {
		return authority, ErrRoleExistence
	}
	copyInfo.Authority.Children = []system.SysAuthority{}
	menus, err := MenuServiceApp.GetMenuAuthority(&request.GetAuthorityId{AuthorityId: copyInfo.OldAuthorityId})
	if err != nil {
		return
	}
	var baseMenu []system.SysBaseMenu
	for _, v := range menus {
		intNum, _ := strconv.Atoi(v.MenuId)
		v.SysBaseMenu.ID = uint(intNum)
		baseMenu = append(baseMenu, v.SysBaseMenu)
	}
	copyInfo.Authority.SysBaseMenus = baseMenu
	err = global.GVA_DB.Create(&copyInfo.Authority).Error
	if err != nil {
		return
	}

	var btns []system.SysAuthorityBtn

	err = global.GVA_DB.Find(&btns, "authority_id = ?", copyInfo.OldAuthorityId).Error
	if err != nil {
		return
	}
	if len(btns) > 0 {
		for i := range btns {
			btns[i].AuthorityId = copyInfo.Authority.AuthorityId
		}
		err = global.GVA_DB.Create(&btns).Error

		if err != nil {
			return
		}
	}
	paths := CasbinServiceApp.GetPolicyPathByAuthorityId(copyInfo.OldAuthorityId)
	err = CasbinServiceApp.UpdateCasbin(copyInfo.Authority.AuthorityId, paths)
	if err != nil {
		_ = authorityService.DeleteAuthority(&copyInfo.Authority)
	}
	return copyInfo.Authority, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: UpdateAuthority
//@description: 更改一个角色
//@param: auth model.SysAuthority
//@return: authority system.SysAuthority, err error

func (authorityService *AuthorityService) UpdateAuthority(auth system.SysAuthority) (authority system.SysAuthority, err error) {
	err = global.GVA_DB.Where("authority_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Updates(&auth).Error
	return auth, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: DeleteAuthority
//@description: 删除角色
//@param: auth *model.SysAuthority
//@return: err error

func (authorityService *AuthorityService) DeleteAuthority(auth *system.SysAuthority) (err error) {
	// 1. 检查角色是否存在
	if errors.Is(global.GVA_DB.Debug().Preload("Users").First(&auth).Error, gorm.ErrRecordNotFound) {
		return errors.New("该角色不存在")
	}
	
	// 2. 检查是否存在子角色
	if !errors.Is(global.GVA_DB.Where("parent_id = ?", auth.AuthorityId).First(&system.SysAuthority{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此角色存在子角色不允许删除")
	}
	
	// 3. 检查 sys_user_authority 多对多关联表
	var userAuthorityCount int64
	if err := global.GVA_DB.Model(&system.SysUserAuthority{}).
		Where("sys_authority_authority_id = ?", auth.AuthorityId).
		Count(&userAuthorityCount).Error; err != nil {
		return errors.New("检查角色关联失败: " + err.Error())
	}
	if userAuthorityCount > 0 {
		return errors.New("此角色有用户正在使用禁止删除（多对多关联表中存在" + strconv.FormatInt(userAuthorityCount, 10) + "条关联记录）")
	}
	
	// 4. 检查 sys_users 表的 authority_id 字段（默认角色）
	var defaultAuthorityCount int64
	if err := global.GVA_DB.Model(&system.SysUser{}).
		Where("authority_id = ?", auth.AuthorityId).
		Count(&defaultAuthorityCount).Error; err != nil {
		return errors.New("检查用户默认角色失败: " + err.Error())
	}
	if defaultAuthorityCount > 0 {
		return errors.New("此角色作为" + strconv.FormatInt(defaultAuthorityCount, 10) + "个用户的默认角色，禁止删除")
	}
	
	// 5. 双重检查：通过 Preload 加载的用户
	if len(auth.Users) != 0 {
		return errors.New("此角色有用户正在使用禁止删除（Preload检测到" + strconv.Itoa(len(auth.Users)) + "个用户）")
	}
	db := global.GVA_DB.Preload("SysBaseMenus").Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(auth)
	err = db.Unscoped().Delete(auth).Error
	if err != nil {
		return
	}
	if len(auth.SysBaseMenus) > 0 {
		err = global.GVA_DB.Model(auth).Association("SysBaseMenus").Delete(auth.SysBaseMenus)
		if err != nil {
			return
		}
		// err = db.Association("SysBaseMenus").Delete(&auth)
	}
	if len(auth.DataAuthorityId) > 0 {
		err = global.GVA_DB.Model(auth).Association("DataAuthorityId").Delete(auth.DataAuthorityId)
		if err != nil {
			return
		}
	}
	err = global.GVA_DB.Delete(&[]system.SysUserAuthority{}, "sys_authority_authority_id = ?", auth.AuthorityId).Error
	if err != nil {
		return
	}
	err = global.GVA_DB.Delete(&[]system.SysAuthorityBtn{}, "authority_id = ?", auth.AuthorityId).Error
	if err != nil {
		return
	}
	authorityId := strconv.Itoa(int(auth.AuthorityId))
	CasbinServiceApp.ClearCasbin(0, authorityId)
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAuthorityInfoList
//@description: 分页获取数据
//@param: info request.PageInfo
//@return: list interface{}, total int64, err error

func (authorityService *AuthorityService) GetAuthorityInfoList(info request.PageInfo) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&system.SysAuthority{})
	if err = db.Where("parent_id = ?", "0").Count(&total).Error; total == 0 || err != nil {
		return
	}
	var authority []system.SysAuthority
	err = db.Limit(limit).Offset(offset).Preload("DataAuthorityId").Where("parent_id = ?", "0").Find(&authority).Error
	for k := range authority {
		err = authorityService.findChildrenAuthority(&authority[k])
	}
	return authority, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAuthorityInfo
//@description: 获取所有角色信息
//@param: auth model.SysAuthority
//@return: sa system.SysAuthority, err error

func (authorityService *AuthorityService) GetAuthorityInfo(auth system.SysAuthority) (sa system.SysAuthority, err error) {
	err = global.GVA_DB.Preload("DataAuthorityId").Where("authority_id = ?", auth.AuthorityId).First(&sa).Error
	return sa, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetDataAuthority
//@description: 设置角色资源权限
//@param: auth model.SysAuthority
//@return: error

func (authorityService *AuthorityService) SetDataAuthority(auth system.SysAuthority) error {
	var s system.SysAuthority
	global.GVA_DB.Preload("DataAuthorityId").First(&s, "authority_id = ?", auth.AuthorityId)
	err := global.GVA_DB.Model(&s).Association("DataAuthorityId").Replace(&auth.DataAuthorityId)
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetMenuAuthority
//@description: 菜单与角色绑定
//@param: auth *model.SysAuthority
//@return: error

func (authorityService *AuthorityService) SetMenuAuthority(auth *system.SysAuthority) error {
	var s system.SysAuthority
	global.GVA_DB.Preload("SysBaseMenus").First(&s, "authority_id = ?", auth.AuthorityId)
	err := global.GVA_DB.Model(&s).Association("SysBaseMenus").Replace(&auth.SysBaseMenus)
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: findChildrenAuthority
//@description: 查询子角色
//@param: authority *model.SysAuthority
//@return: err error

func (authorityService *AuthorityService) findChildrenAuthority(authority *system.SysAuthority) (err error) {
	err = global.GVA_DB.Preload("DataAuthorityId").Where("parent_id = ?", authority.AuthorityId).Find(&authority.Children).Error
	if len(authority.Children) > 0 {
		for k := range authority.Children {
			err = authorityService.findChildrenAuthority(&authority.Children[k])
		}
	}
	return err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetAuthorityUsageInfo
//@description: 获取角色使用情况
//@param: authorityId uint
//@return: map[string]interface{}, error

func (authorityService *AuthorityService) GetAuthorityUsageInfo(authorityId uint) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	
	// 1. 检查角色是否存在
	var authority system.SysAuthority
	if err := global.GVA_DB.Where("authority_id = ?", authorityId).First(&authority).Error; err != nil {
		return nil, errors.New("角色不存在")
	}
	result["authorityInfo"] = authority
	
	// 2. 统计多对多关联表中的用户数量
	var multiUserCount int64
	if err := global.GVA_DB.Model(&system.SysUserAuthority{}).
		Where("sys_authority_authority_id = ?", authorityId).
		Count(&multiUserCount).Error; err != nil {
		return nil, err
	}
	result["multiUserCount"] = multiUserCount
	
	// 3. 统计默认角色为该角色的用户数量
	var defaultUserCount int64
	if err := global.GVA_DB.Model(&system.SysUser{}).
		Where("authority_id = ?", authorityId).
		Where("deleted_at IS NULL").
		Count(&defaultUserCount).Error; err != nil {
		return nil, err
	}
	result["defaultUserCount"] = defaultUserCount
	
	// 4. 获取使用该角色的用户列表（前10个）
	type UserInfo struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		NickName string `json:"nickName"`
		Phone    string `json:"phone"`
	}
	var users []UserInfo
	err := global.GVA_DB.Table("sys_users").
		Select("sys_users.id, sys_users.user_name as username, sys_users.nick_name, sys_users.phone").
		Joins("JOIN sys_user_authority ON sys_users.id = sys_user_authority.sys_user_id").
		Where("sys_user_authority.sys_authority_authority_id = ?", authorityId).
		Where("sys_users.deleted_at IS NULL").
		Limit(10).
		Scan(&users).Error
	if err != nil {
		return nil, err
	}
	result["users"] = users
	
	// 5. 统计子角色数量
	var childrenCount int64
	if err := global.GVA_DB.Model(&system.SysAuthority{}).
		Where("parent_id = ?", authorityId).
		Count(&childrenCount).Error; err != nil {
		return nil, err
	}
	result["childrenCount"] = childrenCount
	
	// 6. 判断是否可以删除
	canDelete := multiUserCount == 0 && defaultUserCount == 0 && childrenCount == 0
	result["canDelete"] = canDelete
	
	// 7. 生成删除阻止原因
	var blockReasons []string
	if multiUserCount > 0 {
		blockReasons = append(blockReasons, strconv.FormatInt(multiUserCount, 10)+"个用户正在使用此角色")
	}
	if defaultUserCount > 0 {
		blockReasons = append(blockReasons, strconv.FormatInt(defaultUserCount, 10)+"个用户以此为默认角色")
	}
	if childrenCount > 0 {
		blockReasons = append(blockReasons, "存在"+strconv.FormatInt(childrenCount, 10)+"个子角色")
	}
	result["blockReasons"] = blockReasons
	
	return result, nil
}
