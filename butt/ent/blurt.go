// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/jeffreycharters/blurt/ent/blurt"
	"github.com/jeffreycharters/blurt/ent/user"
)

// Blurt is the model entity for the Blurt schema.
type Blurt struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreateTime holds the value of the "create_time" field.
	CreateTime time.Time `json:"create_time,omitempty"`
	// UpdateTime holds the value of the "update_time" field.
	UpdateTime time.Time `json:"update_time,omitempty"`
	// Content holds the value of the "content" field.
	Content string `json:"content,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the BlurtQuery when eager-loading is set.
	Edges        BlurtEdges `json:"edges"`
	user_blurts  *uuid.UUID
	selectValues sql.SelectValues
}

// BlurtEdges holds the relations/edges for other nodes in the graph.
type BlurtEdges struct {
	// LikdUsers holds the value of the likd_users edge.
	LikdUsers []*User `json:"likd_users,omitempty"`
	// Author holds the value of the author edge.
	Author *User `json:"author,omitempty"`
	// Liks holds the value of the liks edge.
	Liks []*Lik `json:"liks,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// LikdUsersOrErr returns the LikdUsers value or an error if the edge
// was not loaded in eager-loading.
func (e BlurtEdges) LikdUsersOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.LikdUsers, nil
	}
	return nil, &NotLoadedError{edge: "likd_users"}
}

// AuthorOrErr returns the Author value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e BlurtEdges) AuthorOrErr() (*User, error) {
	if e.Author != nil {
		return e.Author, nil
	} else if e.loadedTypes[1] {
		return nil, &NotFoundError{label: user.Label}
	}
	return nil, &NotLoadedError{edge: "author"}
}

// LiksOrErr returns the Liks value or an error if the edge
// was not loaded in eager-loading.
func (e BlurtEdges) LiksOrErr() ([]*Lik, error) {
	if e.loadedTypes[2] {
		return e.Liks, nil
	}
	return nil, &NotLoadedError{edge: "liks"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Blurt) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case blurt.FieldContent:
			values[i] = new(sql.NullString)
		case blurt.FieldCreateTime, blurt.FieldUpdateTime:
			values[i] = new(sql.NullTime)
		case blurt.FieldID:
			values[i] = new(uuid.UUID)
		case blurt.ForeignKeys[0]: // user_blurts
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Blurt fields.
func (b *Blurt) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case blurt.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				b.ID = *value
			}
		case blurt.FieldCreateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field create_time", values[i])
			} else if value.Valid {
				b.CreateTime = value.Time
			}
		case blurt.FieldUpdateTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field update_time", values[i])
			} else if value.Valid {
				b.UpdateTime = value.Time
			}
		case blurt.FieldContent:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field content", values[i])
			} else if value.Valid {
				b.Content = value.String
			}
		case blurt.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field user_blurts", values[i])
			} else if value.Valid {
				b.user_blurts = new(uuid.UUID)
				*b.user_blurts = *value.S.(*uuid.UUID)
			}
		default:
			b.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Blurt.
// This includes values selected through modifiers, order, etc.
func (b *Blurt) Value(name string) (ent.Value, error) {
	return b.selectValues.Get(name)
}

// QueryLikdUsers queries the "likd_users" edge of the Blurt entity.
func (b *Blurt) QueryLikdUsers() *UserQuery {
	return NewBlurtClient(b.config).QueryLikdUsers(b)
}

// QueryAuthor queries the "author" edge of the Blurt entity.
func (b *Blurt) QueryAuthor() *UserQuery {
	return NewBlurtClient(b.config).QueryAuthor(b)
}

// QueryLiks queries the "liks" edge of the Blurt entity.
func (b *Blurt) QueryLiks() *LikQuery {
	return NewBlurtClient(b.config).QueryLiks(b)
}

// Update returns a builder for updating this Blurt.
// Note that you need to call Blurt.Unwrap() before calling this method if this Blurt
// was returned from a transaction, and the transaction was committed or rolled back.
func (b *Blurt) Update() *BlurtUpdateOne {
	return NewBlurtClient(b.config).UpdateOne(b)
}

// Unwrap unwraps the Blurt entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (b *Blurt) Unwrap() *Blurt {
	_tx, ok := b.config.driver.(*txDriver)
	if !ok {
		panic("ent: Blurt is not a transactional entity")
	}
	b.config.driver = _tx.drv
	return b
}

// String implements the fmt.Stringer.
func (b *Blurt) String() string {
	var builder strings.Builder
	builder.WriteString("Blurt(")
	builder.WriteString(fmt.Sprintf("id=%v, ", b.ID))
	builder.WriteString("create_time=")
	builder.WriteString(b.CreateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("update_time=")
	builder.WriteString(b.UpdateTime.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("content=")
	builder.WriteString(b.Content)
	builder.WriteByte(')')
	return builder.String()
}

// Blurts is a parsable slice of Blurt.
type Blurts []*Blurt