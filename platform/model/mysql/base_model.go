package mysql

import (
	"fmt"
	"github.com/obgnail/plugin-platform/common/log"
	"reflect"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type ModelInter interface {
	New(ModelInter) error
	NewBatch(interface{}) error
	One(ModelInter) error
	All(interface{}, ModelInter) error
	Update(int64, ModelInter) error
	Delete(int64, ModelInter) error
	Save(arg ModelInter) error

	//
	tableName() string
}

type BaseModel struct {
	Id         int64      `gorm:"id" json:"id"`
	CreateTime int64      `gorm:"create_time" json:"create_time"`
	UpdateTime int64      `gorm:"update_time" json:"update_time"`
	Deleted    bool       `gorm:"deleted" json:"deleted"`
	Child      ModelInter `gorm:"-" json:"-"`
}

func (bm *BaseModel) String() string {
	e := reflect.ValueOf(bm.Child).Elem()
	typeOff := e.Type()
	className := typeOff.Name()
	s := "< " + className + "\n"

	var line []string

	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		k := e.Type().Field(i).Name
		v := f.Interface()
		switch f.Kind() {
		case reflect.Struct:
			for j := 0; j < f.NumField(); j++ {
				k := f.Type().Field(j).Name
				v := f.Field(j).Interface()
				if f.Field(j).Kind() == reflect.Interface {
					break
				} else if f.Field(j).Kind() == reflect.Struct {
					break
				}
				l := fmt.Sprintf("%v: %v(%v)", k, reflect.TypeOf(v), v)
				line = append(line, l)
			}
		default:
			l := fmt.Sprintf("%v: %v(%v)", k, reflect.TypeOf(v), v)
			line = append(line, l)
		}
	}

	s = s + strings.Join(line, "\n") + "\n>\n"
	return s
}

func (bm *BaseModel) New(out ModelInter) error {
	className := reflect.ValueOf(out).Elem().Type().Name()

	var s = fmt.Sprintf("%s.Create error: ", className) + "%v"

	now := time.Now().Unix()

	if err := SetAttribute(out, "UpdateTime", now); err != nil {
		return fmt.Errorf(s, err)
	}

	if err := SetAttribute(out, "CreateTime", now); err != nil {
		return fmt.Errorf(s, err)
	}

	if err := SetAttribute(out, "Child", out); err != nil {
		return fmt.Errorf(s, err)
	}

	if err := DB.Create(out).Error; err != nil {
		return fmt.Errorf(s, err)
	}
	return nil
}

func (bm *BaseModel) NewWithDB(db *gorm.DB, out ModelInter) error {
	className := reflect.ValueOf(out).Elem().Type().Name()

	var s = fmt.Sprintf("%s.Create error: ", className) + "%v"

	now := time.Now().Unix()

	if err := SetAttribute(out, "UpdateTime", now); err != nil {
		return fmt.Errorf(s, err)
	}

	if err := SetAttribute(out, "CreateTime", now); err != nil {
		return fmt.Errorf(s, err)
	}

	if err := SetAttribute(out, "Child", out); err != nil {
		return fmt.Errorf(s, err)
	}

	if err := db.Create(out).Error; err != nil {
		return fmt.Errorf(s, err)
	}
	return nil
}

func (bm *BaseModel) One(out ModelInter) error {
	className := reflect.ValueOf(out).Elem().Type().Name()

	var s = fmt.Sprintf("%s.One error: ", className) + "%v"

	err := DB.
		Where("deleted = ?", false).
		Where(out).
		First(out).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
		} else {
			err = fmt.Errorf(s, err)
		}
		return err
	}

	if err := SetAttribute(out, "Child", out); err != nil {
		return fmt.Errorf(s, err)
	}
	return nil
}

func (bm *BaseModel) All(out interface{}, arg ModelInter) error {
	var className = reflect.ValueOf(bm.Child).Elem().Type().Name()

	var db = DB.Table(bm.Child.tableName()).
		Where("deleted = ?", false)
	if arg != nil {
		db = db.Where(arg)
	}
	err := db.Find(out).Error
	if err != nil {
		log.Error("sql: %v", db.SubQuery())
		return fmt.Errorf("%s.All error : %v", className, err)
	}

	e := reflect.ValueOf(out).Elem()
	switch e.Kind() {
	case reflect.Slice:
		for i := 0; i < e.Len(); i++ {
			var o = e.Index(i)
			field := o.Elem().FieldByName("Child")
			field.Set(o)
		}
	}
	return nil
}

func (bm *BaseModel) Delete(id int64, arg ModelInter) error {
	db := DB.Table(bm.Child.tableName()).
		Where("id = ?", id)
	if arg != nil {
		db = db.Where(arg)
	}
	err := db.Update("deleted", true).Error

	if err != nil {
		e := reflect.ValueOf(arg).Elem()
		className := e.Type().Name()
		err = fmt.Errorf("%s.Delete: %v", className, err)
		return err
	}
	return nil
}

func (bm *BaseModel) DeleteWithDB(db *gorm.DB, id int64, arg ModelInter) error {
	orm := db.Table(bm.Child.tableName()).
		Where("id = ?", id)
	if arg != nil {
		orm = orm.Where(arg)
	}
	err := orm.Update("deleted", true).Error

	if err != nil {
		e := reflect.ValueOf(arg).Elem()
		className := e.Type().Name()
		err = fmt.Errorf("%s.Delete: %v", className, err)
		return err
	}
	return nil
}

func (bm *BaseModel) RealDelete(id int64, arg ModelInter) error {
	err := DB.Table(bm.Child.tableName()).
		Where("id = ?", id).Where(arg).
		Delete(bm.Child).
		Error

	if err != nil {
		e := reflect.ValueOf(arg).Elem()
		className := e.Type().Name()
		err = fmt.Errorf("%s.Delete: %v", className, err)
		return err
	}
	return nil
}

func (bm *BaseModel) RealDeleteWithDB(db *gorm.DB, id int64, arg ModelInter) error {
	err := db.Table(bm.Child.tableName()).
		Where("id = ?", id).Where(arg).
		Delete(bm.Child).
		Error

	if err != nil {
		e := reflect.ValueOf(arg).Elem()
		className := e.Type().Name()
		err = fmt.Errorf("%s.Delete: %v", className, err)
		return err
	}
	return nil
}

func (bm *BaseModel) Update(id int64, arg ModelInter) error {
	className := reflect.ValueOf(arg).Elem().Type().Name()

	var s = fmt.Sprintf("%s.Update error: ", className) + "%v"

	//
	now := time.Now().Unix()
	err := SetAttribute(arg, "UpdateTime", now)
	if err != nil {
		return fmt.Errorf(s, err)
	}

	//
	err = DB.Table(bm.Child.tableName()).
		Where("id = ?", id).
		Updates(arg).Error
	if err != nil {
		return fmt.Errorf(s, err)
	}

	//
	if err := bm.One(arg); err != nil {
		if err == gorm.ErrRecordNotFound {
		} else {
			err = fmt.Errorf(s, err)
		}
		return err
	}

	return nil
}

func (bm *BaseModel) UpdateWithDB(db *gorm.DB, id int64, arg ModelInter) error {
	className := reflect.ValueOf(arg).Elem().Type().Name()

	var s = fmt.Sprintf("%s.Update error: ", className) + "%v"

	now := time.Now().Unix()
	err := SetAttribute(arg, "UpdateTime", now)
	if err != nil {
		return fmt.Errorf(s, err)
	}

	err = db.Table(bm.Child.tableName()).
		Where("id = ?", id).
		Updates(arg).Error
	if err != nil {
		return fmt.Errorf(s, err)
	}

	return nil
}

func (bm *BaseModel) Save(arg ModelInter) error {
	err := DB.Save(arg).Error
	if err != nil {
		className := reflect.ValueOf(arg).Elem().Type().Name()
		var s = fmt.Sprintf("%s.Update error: ", className) + "%v"
		return fmt.Errorf(s, err)
	}
	return nil
}

func (bm *BaseModel) NewBatch(data interface{}) error {
	fields, _ := bm.parsedModel(bm.Child)

	var ms = make([]map[string]interface{}, 0)

	e := reflect.Indirect(reflect.ValueOf(data))

	switch e.Kind() {
	case reflect.Slice:
		for i := 0; i < e.Len(); i++ {
			var o = e.Index(i)
			_, m := bm.parsedModel(o.Interface())
			ms = append(ms, m)
		}
	default:
		panic(fmt.Errorf("invalid data type"))
	}

	s := fmt.Sprintf("INSERT INTO %s(", bm.Child.tableName())
	s = s + strings.Join(fields, ", ") + ") VALUES"

	var args []interface{}
	var argPlace []string
	var now = time.Now().Unix()

	for _, m := range ms {
		var l []string

		for _, k := range fields {
			v := m[k]
			if k == "create_time" || k == "update_time" {
				v = now
			}

			args = append(args, v)
			l = append(l, "?")
		}

		line := "(" + strings.Join(l, ", ") + ")"
		argPlace = append(argPlace, line)
	}

	s = s + strings.Join(argPlace, ", ")

	//
	if err := DB.Exec(s, args...).Error; err != nil {
		className := reflect.ValueOf(bm.Child).Elem().Type().Name()
		return fmt.Errorf("%s.NewBatch error : %v", className, err)
	}
	return nil
}

func (bm *BaseModel) NewBatchWithDB(db *gorm.DB, data interface{}) error {
	fields, _ := bm.parsedModel(bm.Child)

	var ms = make([]map[string]interface{}, 0)

	e := reflect.Indirect(reflect.ValueOf(data))

	switch e.Kind() {
	case reflect.Slice:
		for i := 0; i < e.Len(); i++ {
			var o = e.Index(i)
			_, m := bm.parsedModel(o.Interface())
			ms = append(ms, m)
		}
	default:
		panic(fmt.Errorf("invalid data type"))
	}

	s := fmt.Sprintf("INSERT INTO %s(", bm.Child.tableName())
	s = s + strings.Join(fields, ", ") + ") VALUES"

	var args []interface{}
	var argPlace []string
	var now = time.Now().Unix()

	for _, m := range ms {
		var l []string

		for _, k := range fields {
			v := m[k]
			if k == "create_time" || k == "update_time" {
				v = now
			}

			args = append(args, v)
			l = append(l, "?")
		}

		line := "(" + strings.Join(l, ", ") + ")"
		argPlace = append(argPlace, line)
	}

	s = s + strings.Join(argPlace, ", ")
	//
	if err := db.Exec(s, args...).Error; err != nil {
		className := reflect.ValueOf(bm.Child).Elem().Type().Name()
		return fmt.Errorf("%s.NewBatch error : %v", className, err)
	}
	return nil
}

func (bm *BaseModel) parsedModel(data interface{}) ([]string, map[string]interface{}) {
	e := reflect.Indirect(reflect.ValueOf(data))

	fields := make([]string, 0)
	kwargs := make(map[string]interface{})
	for i := 0; i < e.NumField(); i++ {
		f := e.Field(i)
		k := e.Type().Field(i).Tag.Get("gorm")
		if k == "id" {
			continue
		}
		v := f.Interface()

		switch f.Kind() {
		case reflect.Struct:
			ks, ms := bm.parsedModel(v)
			for _, k := range ks {
				v := ms[k]
				fields = append(fields, k)
				kwargs[k] = v
			}
		case reflect.Interface:
			continue
		default:
			fields = append(fields, k)
			kwargs[k] = v
		}
	}
	return fields, kwargs
}

func (bm *BaseModel) Exist(arg ModelInter) (bool, error) {
	err := bm.One(arg)
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}

	if err == gorm.ErrRecordNotFound {
		return false, nil
	}

	return true, nil
}

func GetAttribute(class interface{}, attrName string) (interface{}, error) {
	o := reflect.ValueOf(class)

	if o.Kind() == reflect.Ptr {
		o = o.Elem()
	}

	if o.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid class type")
	}

	r := o.FieldByName(attrName).Interface()
	return r, nil
}

func SetAttribute(class interface{}, attrName string, attrValue interface{}) error {
	o := reflect.ValueOf(class)

	if o.Kind() == reflect.Ptr {
		o = o.Elem()
	}

	if o.Kind() != reflect.Struct {
		return fmt.Errorf("invalid class type")
	}

	f := o.FieldByName(attrName)
	if f.IsValid() && f.CanSet() {
		v := reflect.ValueOf(attrValue)
		f.Set(v)
	} else {
		return fmt.Errorf("invalid attrName")
	}
	return nil
}
