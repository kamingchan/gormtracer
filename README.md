# gormtracer

A [gorm](https://gorm.io/) plugin for opentracing.

_gormtracer only works with GORM 2.0 (`gorm.io/gorm`)_

## Usage

```go
import (
	"github.com/opentracing/opentracing-go"
	"gorm.io/gorm"

	"github.com/kamingchan/gormtracer"
)

func main() {
	// set a tracer
	opentracing.SetGlobalTracer(YourTracer)
	// initialize db session
	db, err := gorm.Open(YourDialector, Config)
	if err != nil {
		// handle error
	}
	// initialize gormtracer
	db.Use(gormtracer.NewGormTracer())

	// happy hacking...
	// db.Where("age > ?", 18).Find()
}
```
