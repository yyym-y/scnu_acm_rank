package model

type Team struct {
	Id         int    `gorm:"column:id;type:int(11);primary_key" form:"id" json:"id"`
	Name       string `gorm:"column:name;type:varchar(255);NOT NULL" form:"name" json:"name"`
	Key        string `gorm:"column:key;type:varchar(255);NOT NULL" form:"key" json:"key"`
	Leader     int64  `gorm:"column:leader;type:bigint(20);NOT NULL" form:"leader" json:"leader"`
	Status     int    `gorm:"column:status;type:tinyint(4);NOT NULL"  json:"status"`
	Ext        string `gorm:"column:ext;type:varchar(255)" json:"ext"`
	NcTeamName string `gorm:"column:nc_team_name;type:varchar(255)" form:"nc_team_name" json:"nc_team_name"`
}

func (m *Team) TableName() string {
	return "team"
}
