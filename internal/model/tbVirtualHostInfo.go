package model

import (
	"time"
)

// TbVirtualHostInfo 主机实例表
type TbVirtualHostInfo struct {
	//mysql.Model `gorm:"embedded"` // embed id and time

	FuniqueID                 uint64     `gorm:"column:funique_id;type:bigint(20);NOT NULL" json:"funiqueId"`                             // 主键
	FinnerInstanceID          string     `gorm:"column:finner_instance_id;type:varchar(50);primary_key" json:"finnerInstanceId"`          // 虚机唯一ID
	FcloudProvider            string     `gorm:"column:fcloud_provider;type:varchar(50);NOT NULL" json:"fcloudProvider"`                  // 服务商标识阿里云：ALI_CLOUD、亚马逊云：AWS_CLOUD、腾讯云：TENCENT_CLOUD、sigma云：SIGMA_CLOUD、非云平台：VIRTUAL_HOST
	Fnation                   string     `gorm:"column:fnation;type:varchar(50)" json:"fnation"`                                          // 归属国家
	Fregion                   string     `gorm:"column:fregion;type:varchar(50)" json:"fregion"`                                          // 归属地区
	Fos                       string     `gorm:"column:fos;type:varchar(50);NOT NULL" json:"fos"`                                         // 操作系统
	FwalIP                    string     `gorm:"column:fwal_ip;type:varchar(15);NOT NULL" json:"fwalIp"`                                  // 公网IP
	FlanIP                    string     `gorm:"column:flan_ip;type:varchar(15);NOT NULL" json:"flanIp"`                                  // 内网IP
	FinstanceName             string     `gorm:"column:finstance_name;type:varchar(100)" json:"finstanceName"`                            // 实例名称
	FinstanceType             string     `gorm:"column:finstance_type;type:varchar(30)" json:"finstanceType"`                             // 实例类型
	FinstanceID               string     `gorm:"column:finstance_id;type:varchar(100)" json:"finstanceId"`                                // 实例ID
	Fcpu                      float32    `gorm:"column:fcpu;type:decimal(20,2);NOT NULL" json:"fcpu"`                                     // cpu 单位: core
	Fmemory                   uint64     `gorm:"column:fmemory;type:bigint(20);NOT NULL" json:"fmemory"`                                  // 内存 单位: byte
	FsystemDisk               uint64     `gorm:"column:fsystem_disk;type:bigint(20);NOT NULL" json:"fsystemDisk"`                         // 系统盘 单位: byte
	FdataDisk                 uint64     `gorm:"column:fdata_disk;type:bigint(20)" json:"fdataDisk"`                                      // 数据盘 单位: byte
	FsshPort                  uint32     `gorm:"column:fssh_port;type:int(11);NOT NULL" json:"fsshPort"`                                  // SSH端口
	FsshUser                  string     `gorm:"column:fssh_user;type:varchar(100);NOT NULL" json:"fsshUser"`                             // SSH登录用户名
	FsshPwd                   string     `gorm:"column:fssh_pwd;type:varchar(100)" json:"fsshPwd"`                                        // ssh用户密码
	FsshPrivate               string     `gorm:"column:fssh_private;type:text" json:"fsshPrivate"`                                        // ssh用户密钥
	FendTime                  *time.Time `gorm:"column:fend_time;type:datetime" json:"fendTime"`                                          // 主机到期时间
	Fstatus                   int        `gorm:"column:fstatus;type:smallint(6);NOT NULL" json:"fstatus"`                                 // 实例状态,取值范围：1-PENDING：表示创建中、2-LAUNCH_FAILED：表示创建失败、3-RUNNING：表示运行中、4-STOPPED：表示关机、5-STARTING：表示开机中、6-STOPPING：表示关机中、7-REBOOTING：表示重启中、8-SHUTTINGDOWN：表示销毁中、9-TERMINATED：表示已销毁
	FlatestOperation          string     `gorm:"column:flatest_operation;type:varchar(50)" json:"flatestOperation"`                       // 创建：CreateInstances、暂停：StopInstances、重启：RebootInstance、销毁：TerminatInstance、启动：RunInstance、扩容磁盘：ResizeInstanceDisks、调整资源：ResetInstancesType、主机信息修改：EditInstance、更改主机密钥：ResetPassword
	FlastestOperationState    int        `gorm:"column:flastest_operation_state;type:smallint(6)" json:"flastestOperationState"`          // 取值范围：1-OPERATING：表示操作执行中、2-SUCCESS：表示操作成功、3-FAILED：表示操作失败
	FlastestOperationTransSeq string     `gorm:"column:flastest_operation_trans_seq;type:varchar(32)" json:"flastestOperationTransSeq"`   // 实例最新操作交易流水号
	FlastestOperationMsg      string     `gorm:"column:flastest_operation_msg;type:text" json:"flastestOperationMsg"`                     // 实例最新操作错误信息
	FcreateTime               time.Time  `gorm:"column:fcreate_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"fcreateTime"` // 创建时间
	FupdateTime               time.Time  `gorm:"column:fupdate_time;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"fupdateTime"` // 修改时间
}

// TableName table name
func (m *TbVirtualHostInfo) TableName() string {
	return "tb_virtual_host_info"
}
