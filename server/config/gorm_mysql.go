package config

import "strings"

type Mysql struct {
	GeneralDB `yaml:",inline" mapstructure:",squash"`
}

func (m *Mysql) Dsn() string {
	cfg := m.Config
	if cfg == "" {
		cfg = "charset=utf8mb4&parseTime=True&loc=Local"
	}
	if !strings.Contains(cfg, "charset=") {
		cfg += "&charset=utf8mb4"
	}
	if !strings.Contains(cfg, "collation=") {
		cfg += "&collation=utf8mb4_unicode_ci"
	}
	return m.Username + ":" + m.Password + "@tcp(" + m.Path + ":" + m.Port + ")/" + m.Dbname + "?" + cfg
}

func (m *Mysql) GetLogMode() string {
	return m.LogMode
}
