package main

import (
	"log"
	"reflect"
	"strings"

	"github.com/go-pg/pg/orm"
	"github.com/stoewer/go-strcase"
)

func GetAllEntities(limit int, offset int, with string, where map[string]string, order string, response interface{}) (err error) {
	if order == "" {
		order = "ASC"
	}

	err = DB().Model(response).
		Apply(Predicate(where)).
		Apply(Pagination(limit, offset)).
		Apply(Relation(with, response)).
		Order("id " + order).
		Select()

	return err
}

func Pagination(limit int, offset int) func(query *orm.Query) (*orm.Query, error) {
	filter := func(query *orm.Query) (*orm.Query, error) {
		if limit == 0 {
			limit = 10
		}
		query = query.Limit(limit)
		if offset != 0 {
			query = query.Offset(offset)
		}
		return query, nil
	}
	return filter
}

func Predicate(where map[string]string) func(query *orm.Query) (*orm.Query, error) {
	if len(where) == 0 {
		return func(query *orm.Query) (*orm.Query, error) {
			return query, nil
		}
	}
	filter := func(query *orm.Query) (*orm.Query, error) {
		for field, value := range where {
			query = query.Where(field+" = ?", value)
		}
		return query, nil
	}
	return filter
}

func GetReflectedType(toReflect interface{}) (reflectedType reflect.Type) {
	reflectedType = reflect.TypeOf(toReflect).Elem()
	if reflectedType.Kind() == reflect.Slice {
		reflectedType = reflectedType.Elem()
	}
	return reflectedType
}

func Relation(with string, response interface{}) func(query *orm.Query) (*orm.Query, error) {

	if with == "" {
		return func(query *orm.Query) (*orm.Query, error) {
			return query, nil
		}
	}
	filter := func(query *orm.Query) (*orm.Query, error) {
		relations := strings.Split(with, ",")
		for _, relation := range relations {
			relation := strcase.UpperCamelCase(relation)

			reflectedType := GetReflectedType(response)

			_, exists := reflectedType.FieldByName(relation)
			if !exists {
				log.Println("field does not exists:", relation)
				continue
			}
			query = query.Relation(relation)
		}
		return query, nil
	}
	return filter
}
