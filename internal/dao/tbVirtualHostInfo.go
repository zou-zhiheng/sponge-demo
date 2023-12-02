package dao

import (
	"context"
	"errors"
	"time"
	"user/internal/cache"
	"user/internal/model"

	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var _ VirtualHostInfoDao = (*virtualHostInfoDao)(nil)

// VirtualHostInfoDao defining the dao interface
type VirtualHostInfoDao interface {
	Create(ctx context.Context, table *model.TbVirtualHostInfo) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, table *model.TbVirtualHostInfo) error
	GetByID(ctx context.Context, innerInstanceId string) (*model.TbVirtualHostInfo, error)
	GetByCondition(ctx context.Context, condition *query.Conditions) (*model.TbVirtualHostInfo, error)
	GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.TbVirtualHostInfo, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.TbVirtualHostInfo, int64, error)

	CreateByTx(ctx context.Context, tx *gorm.DB, table *model.TbVirtualHostInfo) (uint64, error)
	UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.TbVirtualHostInfo) error
	DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error
	ExtendVirtualHostInfoDao
}

type ExtendVirtualHostInfoDao interface {
	BatchCreate(ctx context.Context, tables []*model.TbVirtualHostInfo) error
	GetByTransSeq(ctx context.Context, transSeq string) ([]*model.TbVirtualHostInfo, int64, error)
}

type virtualHostInfoDao struct {
	db    *gorm.DB
	cache cache.TbVirtualHostInfoCache
	sfg   *singleflight.Group
}

func (d *virtualHostInfoDao) BatchCreate(ctx context.Context, tables []*model.TbVirtualHostInfo) error {
	for _, t := range tables {
		err := d.db.WithContext(ctx).Create(t).Error
		_ = d.cache.Del(ctx, t.FuniqueID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *virtualHostInfoDao) GetByTransSeq(ctx context.Context, transSeq string) ([]*model.TbVirtualHostInfo, int64, error) {
	//TODO implement me
	panic("implement me")
}

// NewTbVirtualHostInfoDao creating the dao interface
func NewTbVirtualHostInfoDao(db *gorm.DB) VirtualHostInfoDao {
	return &virtualHostInfoDao{
		db: db,
		//cache: xCache,
		sfg: new(singleflight.Group),
	}
}

// Create a record, insert the record and the id value is written back to the table
func (d *virtualHostInfoDao) Create(ctx context.Context, table *model.TbVirtualHostInfo) error {
	err := d.db.WithContext(ctx).Create(table).Error
	//_ = d.cache.Del(ctx, table.FuniqueID)
	return err
}

// DeleteByID delete a record by id
func (d *virtualHostInfoDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.TbVirtualHostInfo{}).Error
	if err != nil {
		return err
	}

	// delete cache
	//_ = d.cache.Del(ctx, id)

	return nil
}

// DeleteByIDs delete records by batch id
func (d *virtualHostInfoDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.TbVirtualHostInfo{}).Error
	if err != nil {
		return err
	}

	// delete cache
	//for _, id := range ids {
	//	_ = d.cache.Del(ctx, id)
	//}

	return nil
}

// UpdateByID update a record by id
func (d *virtualHostInfoDao) UpdateByID(ctx context.Context, table *model.TbVirtualHostInfo) error {
	err := d.updateDataByID(ctx, d.db, table)

	// delete cache
	//_ = d.cache.Del(ctx, table.ID)

	return err
}

func (d *virtualHostInfoDao) updateDataByID(ctx context.Context, db *gorm.DB, table *model.TbVirtualHostInfo) error {
	if table.FuniqueID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.FuniqueID != 0 {
		update["funique_id"] = table.FuniqueID
	}
	if table.FinnerInstanceID != "" {
		update["finner_instance_id"] = table.FinnerInstanceID
	}
	if table.FcloudProvider != "" {
		update["fcloud_provider"] = table.FcloudProvider
	}
	if table.Fnation != "" {
		update["fnation"] = table.Fnation
	}
	if table.Fregion != "" {
		update["fregion"] = table.Fregion
	}
	if table.Fos != "" {
		update["fos"] = table.Fos
	}
	if table.FwalIP != "" {
		update["fwal_ip"] = table.FwalIP
	}
	if table.FlanIP != "" {
		update["flan_ip"] = table.FlanIP
	}
	if table.FinstanceName != "" {
		update["finstance_name"] = table.FinstanceName
	}
	if table.FinstanceType != "" {
		update["finstance_type"] = table.FinstanceType
	}
	if table.FinstanceID != "" {
		update["finstance_id"] = table.FinstanceID
	}
	if table.Fcpu != 0 {
		update["fcpu"] = table.Fcpu
	}
	if table.Fmemory != 0 {
		update["fmemory"] = table.Fmemory
	}
	if table.FsystemDisk != 0 {
		update["fsystem_disk"] = table.FsystemDisk
	}
	if table.FdataDisk != 0 {
		update["fdata_disk"] = table.FdataDisk
	}
	if table.FsshPort != 0 {
		update["fssh_port"] = table.FsshPort
	}
	if table.FsshUser != "" {
		update["fssh_user"] = table.FsshUser
	}
	if table.FsshPwd != "" {
		update["fssh_pwd"] = table.FsshPwd
	}
	if table.FsshPrivate != "" {
		update["fssh_private"] = table.FsshPrivate
	}
	if table.FendTime != nil && table.FendTime.IsZero() == false {
		update["fend_time"] = table.FendTime
	}
	if table.Fstatus != 0 {
		update["fstatus"] = table.Fstatus
	}
	if table.FlatestOperation != "" {
		update["flatest_operation"] = table.FlatestOperation
	}
	if table.FlastestOperationState != 0 {
		update["flastest_operation_state"] = table.FlastestOperationState
	}
	if table.FlastestOperationTransSeq != "" {
		update["flastest_operation_trans_seq"] = table.FlastestOperationTransSeq
	}
	if table.FlastestOperationMsg != "" {
		update["flastest_operation_msg"] = table.FlastestOperationMsg
	}
	//if table.FcreateTime.IsZero() == false {
	//	update["fcreate_time"] = table.FcreateTime
	//}
	//if table.FupdateTime.IsZero() == false {
	//	update["fupdate_time"] = table.FupdateTime
	//}

	return db.WithContext(ctx).Model(table).Updates(update).Error
}

// GetByID get a record by id
func (d *virtualHostInfoDao) GetByID(ctx context.Context, id string) (*model.TbVirtualHostInfo, error) {
	// for the same id, prevent high concurrent simultaneous access to mysql
	val, err, _ := d.sfg.Do(id, func() (interface{}, error) { //nolint
		table := &model.TbVirtualHostInfo{}
		err := d.db.WithContext(ctx).Where("finner_instance_id = ?", id).First(table).Error
		if err != nil {
			// if data is empty, set not found cache to prevent cache penetration, default expiration time 10 minutes
			if errors.Is(err, model.ErrRecordNotFound) {
				//err = d.cache.SetCacheWithNotFound(ctx, id)
				//if err != nil {
				//	return nil, err
				//}
				return nil, model.ErrRecordNotFound
			}
			return nil, err
		}
		// set cache
		//err = d.cache.Set(ctx, id, table, cache.TbVirtualHostInfoExpireTime)
		//if err != nil {
		//	return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
		//}
		return table, nil
	})
	if err != nil {
		return nil, err
	}
	table, ok := val.(*model.TbVirtualHostInfo)
	if !ok {
		return nil, model.ErrRecordNotFound
	}
	return table, nil
}

// GetByCondition get a record by condition
// query conditions:
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like
//	value: column name
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: find a male aged 20
//
//	condition = &query.Conditions{
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *virtualHostInfoDao) GetByCondition(ctx context.Context, c *query.Conditions) (*model.TbVirtualHostInfo, error) {
	queryStr, args, err := c.ConvertToGorm()
	if err != nil {
		return nil, err
	}

	table := &model.TbVirtualHostInfo{}
	err = d.db.WithContext(ctx).Where(queryStr, args...).First(table).Error
	if err != nil {
		return nil, err
	}

	return table, nil
}

// GetByIDs list of records by batch id
func (d *virtualHostInfoDao) GetByIDs(ctx context.Context, ids []uint64) (map[uint64]*model.TbVirtualHostInfo, error) {
	//itemMap, err := d.cache.MultiGet(ctx, ids)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var missedIDs []uint64
	//for _, id := range ids {
	//	_, ok := itemMap[id]
	//	if !ok {
	//		missedIDs = append(missedIDs, id)
	//		continue
	//	}
	//}
	//
	//// get missed data
	//if len(missedIDs) > 0 {
	//	// find the id of an active placeholder, i.e. an id that does not exist in mysql
	//	var realMissedIDs []uint64
	//	for _, id := range missedIDs {
	//		_, err = d.cache.Get(ctx, id)
	//		if errors.Is(err, cacheBase.ErrPlaceholder) {
	//			continue
	//		}
	//		realMissedIDs = append(realMissedIDs, id)
	//	}
	//
	//	if len(realMissedIDs) > 0 {
	//		var missedData []*model.TbVirtualHostInfo
	//		err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
	//		if err != nil {
	//			return nil, err
	//		}
	//
	//		if len(missedData) > 0 {
	//			//for _, data := range missedData {
	//			//	itemMap[data.ID] = data
	//			//}
	//			err = d.cache.MultiSet(ctx, missedData, cache.TbVirtualHostActionInfoExpireTime)
	//			if err != nil {
	//				return nil, err
	//			}
	//		} else {
	//			for _, id := range realMissedIDs {
	//				_ = d.cache.SetCacheWithNotFound(ctx, id)
	//			}
	//		}
	//	}
	//}

	//return itemMap, nil
	return nil, nil
}

// GetByColumns get records by paging and column information,
// Note: suitable for scenarios where the number of rows in the table is not very large,
//
//	performance is lower if the data table is large because of the use of offset.
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	size: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default is "=",  support =, !=, >, >=, <, <=, like
//	value: column name
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Size: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *virtualHostInfoDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.TbVirtualHostInfo, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Unscoped().Model(&model.TbVirtualHostInfo{}).Select([]string{"funique_id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.TbVirtualHostInfo{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

// CreateByTx create a record in the database using the provided transaction
func (d *virtualHostInfoDao) CreateByTx(ctx context.Context, tx *gorm.DB, table *model.TbVirtualHostInfo) (uint64, error) {
	err := tx.WithContext(ctx).Create(table).Error
	return table.FuniqueID, err
}

// DeleteByTx delete a record by id in the database using the provided transaction
func (d *virtualHostInfoDao) DeleteByTx(ctx context.Context, tx *gorm.DB, id uint64) error {
	update := map[string]interface{}{
		"deleted_at": time.Now(),
	}
	err := tx.WithContext(ctx).Model(&model.TbVirtualHostInfo{}).Where("id = ?", id).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	//_ = d.cache.Del(ctx, id)

	return nil
}

// UpdateByTx update a record by id in the database using the provided transaction
func (d *virtualHostInfoDao) UpdateByTx(ctx context.Context, tx *gorm.DB, table *model.TbVirtualHostInfo) error {
	err := d.updateDataByID(ctx, tx, table)

	// delete cache
	//_ = d.cache.Del(ctx, table.ID)

	return err
}
