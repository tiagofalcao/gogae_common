package engine

import (
	"appengine"
	"appengine/datastore"

	"errors"
	"reflect"
)

func EntityPrint(c appengine.Context, p interface{}) {
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		c.Infof("Invalid Type")
	}
	v = v.Elem()
	t := v.Type()
	c.Infof("Type: %s", t)
	c.Infof("Setable: %s", v.CanSet())
	for i := 0; i < t.NumField(); i++ {
		f := v.Field(i)
		c.Infof("%d: %s %s = %v", i, t.Field(i).Name, f.Type(), f.Interface())
	}
}

func GetKey(p interface{}) (key *datastore.Key) { // TODO add error
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return
	}
	v = v.Elem()
	if !v.CanSet() {
		return
	}
	k := v.FieldByName("Key")
	i := k.Interface()
	switch r := i.(type) {
	case *datastore.Key:
		{
			return r
		}
	}
	return nil
}

func SetKey(p interface{}, key *datastore.Key) { //TODO add error
	v := reflect.ValueOf(p)
	if v.Kind() != reflect.Ptr || v.IsNil() || v.Elem().Kind() != reflect.Struct {
		return
	}
	v = v.Elem()
	if !v.CanSet() {
		return
	}
	kp := v.FieldByName("Key")
	k := reflect.ValueOf(key)
	kp.Set(k)
}

func Put(c appengine.Context, key *datastore.Key, src interface{}) error {
	fullkey, err := datastore.Put(c, key, src)
	if err == nil {
		SetKey(src, fullkey)
	}
	return err
}

func Save(c appengine.Context, src interface{}) error {
	key := GetKey(src)
	if key == nil {
		return errors.New("No key found to save")
	}
	_, err := datastore.Put(c, key, src)
	return err
}

func Get(c appengine.Context, key *datastore.Key, dst interface{}) error {
	err := datastore.Get(c, key, dst)
	if err != nil {
		SetKey(&dst, key)
	}
	return err

}
