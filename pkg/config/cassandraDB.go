package config

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

type Cassandra struct {
	Session *gocql.Session
}

func SetupCassandra() (*Cassandra, error) {
	var cassandra *Cassandra = &Cassandra{}
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "mantritva"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Unable to connect cassandra: %v", err)
		return nil, err
	}
	cassandra.Session = session
	if err = cassandra.SetupTables(); err != nil {
		return nil, err
	}
	return cassandra, nil
}

func (c *Cassandra) Close() {
	if c.Session != nil {
		c.Session.Close()
	}
}

func (c *Cassandra) Ping() error {
	if c.Session == nil {
		return errors.New("Cassandra session is not initialized")
	}
	var ts gocql.UUID
	if err := c.Session.Query("SELECT now() FROM system.local").Scan(&ts); err != nil {
		return err
	}
	fmt.Println(ts)
	return nil
}

func (c *Cassandra) SetupTables() error {
	var err error
	if err = c.CreateNewsPressPostTable(); err != nil {
		log.Println(err)
		return err
	}
	if err = c.CreatePostReactionsTable(); err != nil {
		log.Println(err)
		return err

	}
	if err = c.CreateFeedsTable(); err != nil {
		log.Println(err)
		return err
	}
	// if err = c.CreateIndices(); err != nil {
	// 	log.Println(err)
	// 	return err
	// }
	return nil
}

func (c *Cassandra) Save(table string, data interface{}) error {
	query, values := PrepareInsertQuery(table, data)
	return c.Session.Query(query, values...).Exec()
}

func (c *Cassandra) GetRows(query string, out interface{}) error {
	iter := c.Session.Query(query).Iter()
	row := make(map[string]interface{})
	sliceValue := reflect.ValueOf(out).Elem()
	elemType := sliceValue.Type().Elem()
	for iter.MapScan(row) {
		elemPtr := reflect.New(elemType)
		mapRowToStruct(row, elemPtr.Interface())
		sliceValue.Set(reflect.Append(sliceValue, elemPtr.Elem()))
		row = make(map[string]interface{})
	}
	return iter.Close()
}

func (c *Cassandra) GetRow(query string, out interface{}) (bool, error) {
	row := make(map[string]interface{})
	if err := c.Session.Query(query).Consistency(gocql.One).MapScan(row); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, err
	}
	mapRowToStruct(row, out)
	return true, nil
}

func (c *Cassandra) CreateIndices() error {
	var err error
	newsPressPostsIndexQuery := `CREATE INDEX news_press_posts_id_memberId_partyId_idx ON mantritva.news_press_posts (id, memberId, partyId);`
	if err = c.Session.Query(newsPressPostsIndexQuery).Exec(); err != nil {
		return err
	}
	postReactionsIndexQuery := `CREATE INDEX post_reactions_id_partyId_partyId_idx ON mantritva.post_reactions (id, partyId, userId);`
	if err = c.Session.Query(postReactionsIndexQuery).Exec(); err != nil {
		return err
	}
	return nil
}

func (c *Cassandra) CreateFeedsTable() error {
	query := `CREATE TABLE IF NOT EXISTS mantritva.feeds(
				id TIMEUUID PRIMARY KEY, 
				type TEXT, 
				postId TIMEUUID,
				createdBy TEXT,
				createdAt TIMESTAMP, 
				updatedAt TIMESTAMP,
				isDeleted BOOLEAN, 
				deletedAt TIMESTAMP);`
	return c.Session.Query(query).Exec()
}

func (c *Cassandra) CreateNewsPressPostTable() error {
	query := `CREATE TABLE IF NOT EXISTS mantritva.news_press_posts(
				id TIMEUUID PRIMARY KEY, 
				memberId int, 
				partyId int, 
				subject TEXT, 
				type TEXT, 
				tone TEXT, 
				statement TEXT, 
				statementDate date,
				status TEXT,
				upvoteCount int,
				downvoteCount int,
				commentCount int, 
				createdBy TEXT,
				createdAt TIMESTAMP, 
				updatedAt TIMESTAMP,
				isDeleted BOOLEAN, 
				deletedAt TIMESTAMP);`
	return c.Session.Query(query).Exec()
}

func (c *Cassandra) CreatePostReactionsTable() error {
	query := `CREATE TABLE IF NOT EXISTS mantritva.post_reactions(
				id TIMEUUID, 
				postId TIMEUUID, 
				userId int, 
				reaction TEXT,
				userIP TEXT,
				createdBy TEXT,
				createdAt TIMESTAMP, 
				updatedAt TIMESTAMP,
				isDeleted BOOLEAN, 
				deletedAt TIMESTAMP,
				PRIMARY KEY ((postid), userid));`
	return c.Session.Query(query).Exec()
}

func extractFields(v reflect.Value, t reflect.Type, columns *[]string, placeholders *[]string, values *[]interface{}) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		valField := v.Field(i)
		// Skip unexported fields
		if field.PkgPath != "" {
			continue
		}
		if field.Anonymous {
			ft := field.Type
			fv := valField
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				if fv.IsNil() {
					continue
				}
				fv = fv.Elem()
			}
			extractFields(fv, ft, columns, placeholders, values)
			continue
		}
		col := field.Tag.Get("json")
		tagName := ""
		if col != "" {
			tagName = strings.Split(col, ",")[0]
		}
		if tagName != "" && tagName != "-" {
			tagName = strings.ReplaceAll(tagName, "_", "")
			*columns = append(*columns, tagName)
			*placeholders = append(*placeholders, "?")
			if field.Type.String() != "gorm.DeletedAt" {
				if tagName == "createdat" {
					*values = append(*values, time.Now())
				} else {
					*values = append(*values, v.Field(i).Interface())
				}

			} else {
				*values = append(*values, nil)
			}

		}
	}
}

func PrepareInsertQuery(table string, data interface{}) (string, []interface{}) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	var columns []string
	var placeholders []string
	var values []interface{}
	extractFields(v, t, &columns, &placeholders, &values)
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);",
		table,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
	)
	return query, values
}

func mapRowToStruct(row map[string]interface{}, out interface{}) {
	v := reflect.ValueOf(out).Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := v.Field(i)
		// Handle anonymous (embedded) fields recursively
		if field.Anonymous {
			// Dereference pointer if needed
			ft := field.Type
			fvv := fv
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
				if fvv.IsNil() {
					fvv.Set(reflect.New(ft))
				}
				fvv = fvv.Elem()
			}
			mapRowToStruct(row, fvv.Addr().Interface())
			continue
		}
		// Use json tag if present, else field name
		tag := field.Tag.Get("json")
		col := tag
		if idx := strings.Index(tag, ","); idx != -1 {
			col = tag[:idx]
		}
		if col == "" {
			col = strings.ToLower(field.Name)
		} else {
			col = strings.ToLower(col)
		}
		col = strings.ReplaceAll(col, "_", "")
		if val, ok := row[col]; ok && val != nil {
			if fv.CanSet() {
				if fv.Type() == reflect.TypeOf(time.Time{}) {
					switch v := val.(type) {
					case time.Time:
						fv.Set(reflect.ValueOf(v))
					case *time.Time:
						if v != nil {
							fv.Set(reflect.ValueOf(*v))
						}
					case string:
						parsed, err := time.Parse(time.RFC3339, v)
						if err == nil {
							fv.Set(reflect.ValueOf(parsed))
						}
					}
					continue
				}
				// Handle type conversion if needed
				valValue := reflect.ValueOf(val)
				if valValue.Type().ConvertibleTo(fv.Type()) {
					fv.Set(valValue.Convert(fv.Type()))
				}
			}
		}
	}
}
