package supervillain

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFieldName(t *testing.T) {
	assert.Equal(t,
		fieldName(reflect.StructField{Name: "RCONPassword"}),
		"RCONPassword",
	)

	assert.Equal(t,
		fieldName(reflect.StructField{Name: "LANMode"}),
		"LANMode",
	)

	assert.Equal(t,
		fieldName(reflect.StructField{Name: "ABC"}),
		"ABC",
	)
}

func TestFieldNameJsonTag(t *testing.T) {
	type S struct {
		NotTheFieldName string `json:"fieldName"`
	}

	assert.Equal(t,
		fieldName(reflect.TypeOf(S{}).Field((0))),
		"fieldName",
	)
}

func TestFieldNameJsonTagOmitEmpty(t *testing.T) {
	type S struct {
		NotTheFieldName string `json:"fieldName,omitempty"`
	}

	assert.Equal(t,
		fieldName(reflect.TypeOf(S{}).Field((0))),
		"fieldName",
	)
}

func TestSchemaName(t *testing.T) {
	assert.Equal(t,
		schemaName("", "User"),
		"UserSchema",
	)
	assert.Equal(t,
		schemaName("Bot", "User"),
		"BotUserSchema",
	)
}

func TestStructSimple(t *testing.T) {
	type User struct {
		Name   string
		Age    int
		Height float64
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Age: z.number(),
  Height: z.number(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStructSimpleWithOmittedField(t *testing.T) {
	type User struct {
		Name        string
		Age         int
		Height      float64
		NotExported string `json:"-"`
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Age: z.number(),
  Height: z.number(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStructSimplePrefix(t *testing.T) {
	type User struct {
		Name   string
		Age    int
		Height float64
	}
	assert.Equal(t,
		`export const BotUserSchema = z.object({
  Name: z.string(),
  Age: z.number(),
  Height: z.number(),
})
export type BotUser = z.infer<typeof BotUserSchema>

`,
		StructToZodSchemaWithPrefix("Bot", User{}))
}

func TestStringArray(t *testing.T) {
	type User struct {
		Tags []string
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Tags: z.string().array().nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStructArray(t *testing.T) {
	type User struct {
		Favourites []struct {
			Name string
		}
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Favourites: z.object({
    Name: z.string(),
  }).array().nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStructArrayOptional(t *testing.T) {
	type User struct {
		Favourites []struct {
			Name string
		} `json:",omitempty"`
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Favourites: z.object({
    Name: z.string(),
  }).array().optional(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStructArrayOptionalNullable(t *testing.T) {
	type User struct {
		Favourites *[]struct {
			Name string
		} `json:",omitempty"`
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Favourites: z.object({
    Name: z.string(),
  }).array().optional().nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStringOptional(t *testing.T) {
	type User struct {
		Name     string
		Nickname string `json:",omitempty"`
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Nickname: z.string().optional(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStringNullable(t *testing.T) {
	type User struct {
		Name     string
		Nickname *string
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Nickname: z.string().nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStringOptionalNotNullable(t *testing.T) {
	type User struct {
		Name     string
		Nickname *string `json:",omitempty"` // nil values are omitted
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Nickname: z.string().optional(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStringOptionalNullable(t *testing.T) {
	type User struct {
		Name     string
		Nickname **string `json:",omitempty"` // nil values are omitted
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Nickname: z.string().optional().nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestStringArrayNullable(t *testing.T) {
	type User struct {
		Name string
		Tags []*string
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Tags: z.string().array().nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestInterfaceAny(t *testing.T) {
	type User struct {
		Name     string
		Metadata interface{}
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Metadata: z.any(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestInterfacePointerAny(t *testing.T) {
	type User struct {
		Name     string
		Metadata *interface{}
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Metadata: z.any(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestInterfaceEmptyAny(t *testing.T) {
	type User struct {
		Name     string
		Metadata interface{} `json:",omitempty"`
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Metadata: z.any(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestInterfacePointerEmptyAny(t *testing.T) {
	type User struct {
		Name     string
		Metadata *interface{} `json:",omitempty"`
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Metadata: z.any(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestMapStringToString(t *testing.T) {
	type User struct {
		Name     string
		Metadata map[string]string
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Metadata: z.record(z.string(), z.string()).nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestMapStringToInterface(t *testing.T) {
	type User struct {
		Name     string
		Metadata map[string]interface{}
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Metadata: z.record(z.string(), z.any()).nullable(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestMapWithStruct(t *testing.T) {

	type PostWithMetaData struct {
		Title string
	}
	type User struct {
		MapWithStruct map[string]PostWithMetaData
	}
	assert.Equal(t,
		`export const PostWithMetaDataSchema = z.object({
  Title: z.string(),
})
export type PostWithMetaData = z.infer<typeof PostWithMetaDataSchema>

export const UserSchema = z.object({
  MapWithStruct: z.record(z.string(), PostWithMetaDataSchema).nullable(),
})
export type User = z.infer<typeof UserSchema>

`, StructToZodSchema(User{}))
}

func TestEverything(t *testing.T) {
	// The order matters PostWithMetaData needs to be declared after post otherwise it will raise a
	// `Block-scoped variable 'Post' used before its declaration.` typescript error.
	type Post struct {
		Title string
	}
	type PostWithMetaData struct {
		Title string
		Post  Post
	}
	type User struct {
		Name                 string
		Nickname             *string // pointers become optional
		Age                  int
		Height               float64
		OldPostWithMetaData  PostWithMetaData
		Tags                 []string
		TagsOptional         []string   `json:",omitempty"` // slices with omitempty cannot be null
		TagsOptionalNullable *[]string  `json:",omitempty"` // pointers to slices with omitempty can be null or undefined
		Favourites           []struct { // nested structs are kept inline
			Name string
		}
		Posts                         []Post             // external structs are emitted as separate exports
		Post                          Post               `json:",omitempty"` // this tag is ignored because structs don't have an empty value
		PostOptional                  *Post              `json:",omitempty"` // single struct pointers with omitempty cannot be null
		PostOptionalNullable          **Post             `json:",omitempty"` // double struct pointers with omitempty can be null
		Metadata                      map[string]string  // maps can be null
		MetadataOptional              map[string]string  `json:",omitempty"` // maps with omitempty cannot be null
		MetadataOptionalNullable      *map[string]string `json:",omitempty"` // pointers to maps with omitempty can be null or undefined
		ExtendedProps                 interface{}        // interfaces are just "any" even though they can be null
		ExtendedPropsOptional         interface{}        `json:",omitempty"` // interfaces with omitempty are still just "any"
		ExtendedPropsNullable         *interface{}       // pointers to interfaces are just "any"
		ExtendedPropsOptionalNullable *interface{}       `json:",omitempty"` // pointers to interfaces with omitempty are also just "any"
		ExtendedPropsVeryIndirect     ****interface{}    // interfaces are always "any" no matter the levels of indirection
		NewPostWithMetaData           PostWithMetaData
		VeryNewPost                   Post
		MapWithStruct                 map[string]PostWithMetaData
	}
	assert.Equal(t,
		`export const PostSchema = z.object({
  Title: z.string(),
})
export type Post = z.infer<typeof PostSchema>

export const PostWithMetaDataSchema = z.object({
  Title: z.string(),
  Post: PostSchema,
})
export type PostWithMetaData = z.infer<typeof PostWithMetaDataSchema>

export const UserSchema = z.object({
  Name: z.string(),
  Nickname: z.string().nullable(),
  Age: z.number(),
  Height: z.number(),
  OldPostWithMetaData: PostWithMetaDataSchema,
  Tags: z.string().array().nullable(),
  TagsOptional: z.string().array().optional(),
  TagsOptionalNullable: z.string().array().optional().nullable(),
  Favourites: z.object({
    Name: z.string(),
  }).array().nullable(),
  Posts: PostSchema.array().nullable(),
  Post: PostSchema,
  PostOptional: PostSchema.optional(),
  PostOptionalNullable: PostSchema.optional().nullable(),
  Metadata: z.record(z.string(), z.string()).nullable(),
  MetadataOptional: z.record(z.string(), z.string()).optional(),
  MetadataOptionalNullable: z.record(z.string(), z.string()).optional().nullable(),
  ExtendedProps: z.any(),
  ExtendedPropsOptional: z.any(),
  ExtendedPropsNullable: z.any(),
  ExtendedPropsOptionalNullable: z.any(),
  ExtendedPropsVeryIndirect: z.any(),
  NewPostWithMetaData: PostWithMetaDataSchema,
  VeryNewPost: PostSchema,
  MapWithStruct: z.record(z.string(), PostWithMetaDataSchema).nullable(),
})
export type User = z.infer<typeof UserSchema>

`, StructToZodSchema(User{}))
}

func TestConvertSlice(t *testing.T) {
	type Foo struct {
		Bar string
		Baz string
		Quz string
	}

	type Zip struct {
		Zap *Foo
	}

	type Whim struct {
		Wham *Foo
	}
	c := NewConverter(map[string]CustomFn{})
	types := []interface{}{
		Zip{},
		Whim{},
	}
	assert.Equal(t,
		`export const ZipSchema = z.object({
  Zap: FooSchema.nullable(),
})
export type Zip = z.infer<typeof ZipSchema>

export const FooSchema = z.object({
  Bar: z.string(),
  Baz: z.string(),
  Quz: z.string(),
})
export type Foo = z.infer<typeof FooSchema>

export const WhimSchema = z.object({
  Wham: FooSchema.nullable(),
})
export type Whim = z.infer<typeof WhimSchema>

`, c.ConvertSlice(types))
}

func TestStructTime(t *testing.T) {
	type User struct {
		Name string
		When time.Time
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  When: z.string(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

func TestCustom(t *testing.T) {
	c := NewConverter(map[string]CustomFn{
		"github.com/Southclaws/supervillain.Decimal": func(c *Converter, t reflect.Type, s, g string, i int) string {
			return "z.string()"
		},
	})

	type Decimal struct {
		value    int
		exponent int
	}

	type User struct {
		Name  string
		Money Decimal
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Name: z.string(),
  Money: z.string(),
})
export type User = z.infer<typeof UserSchema>

`,
		c.Convert(User{}))
}

type State int

func (s State) ZodSchema() string {
	return "z.string()"
}

func TestZodSchemaConstant(t *testing.T) {
	type Job struct {
		State State
	}
	assert.Equal(t,
		`export const JobSchema = z.object({
  State: z.string(),
})
export type Job = z.infer<typeof JobSchema>

`,
		StructToZodSchema(Job{}))
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) ZodSchema(c *Converter, t reflect.Type, name, generic string, indent int) string {
	return fmt.Sprintf("%s.array()", c.ConvertType(t.Key(), name, indent))
}

func TestZodSchemaDynamic(t *testing.T) {
	type User struct {
		Nicknames       Set[string]
		FavoriteNumbers Set[int]
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Nicknames: z.string().array(),
  FavoriteNumbers: z.number().array(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

type Set2[T comparable] map[T]struct{}

func (s Set2[T]) ZodSchema(convert func(t reflect.Type, name string, indent int) string, t reflect.Type, name, generic string, indent int) string {
	return fmt.Sprintf("%s.array()", convert(t.Key(), name, indent))
}

func TestZodSchemaDynamicFunction(t *testing.T) {
	type User struct {
		Nicknames       Set2[string]
		FavoriteNumbers Set2[int]
	}
	assert.Equal(t,
		`export const UserSchema = z.object({
  Nicknames: z.string().array(),
  FavoriteNumbers: z.number().array(),
})
export type User = z.infer<typeof UserSchema>

`,
		StructToZodSchema(User{}))
}

type Strange int

func (s Strange) ZodSchema(weird int, signature string) int {
	return int(s)
}

func TestZodSchemaUnexpected(t *testing.T) {
	type User struct {
		Strange Strange
	}
	assert.Panics(t, func() {
		StructToZodSchema(User{})
	})
}

type StateWithoutSchema int

func (s StateWithoutSchema) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s StateWithoutSchema) String() string {
	return fmt.Sprint(int(s))
}

func TestStrictCustom(t *testing.T) {
	type Job struct {
		State StateWithoutSchema
	}

	c := NewConverter(map[string]CustomFn{}, WithStrictCustomSchemas(true))
	assert.Panics(t, func() {
		c.Convert(Job{})
	})

	c2 := NewConverter(map[string]CustomFn{}, WithStrictCustomSchemas(false))
	assert.Equal(t,
		`export const JobSchema = z.object({
  State: z.number(),
})
export type Job = z.infer<typeof JobSchema>

`, c2.Convert(Job{}))

	c3 := NewConverter(map[string]CustomFn{} /* defaults to false */)
	assert.Equal(t,
		`export const JobSchema = z.object({
  State: z.number(),
})
export type Job = z.infer<typeof JobSchema>

`, c3.Convert(Job{}))
}
