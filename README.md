### How to create a new Service
* Create `{service_name}.go` file
* Open `{service_name}.go` file
    - Create `type {service_name} interface{}`
* Open `setup.go` file
    - Create `type {service_name}Impl struct{}`
    - Create `type {service_name}Opts struct{}`
    - Create `func New{service_name} {service_name}` 
* Open `init.go` file
    - Add your `{service_name}` in `ServiceImpl` 
    - Create a function `Get{service_name}() {service_name}` inside `Service` interface.
    - Implement `Get{service_name}() {service_name}` for `ServiceImpl`
    - Add `si.{service_name}` inside `setup(opts *ServiceOpts)` to initialize the new service for `Service` package using `func New{service_name} {service_name}`
* Open `{service_name}.go` file 
    - Start adding functions in `type {service_name} interface` and implement using `type {service_name}Impl struct`

# Unit Test
## Service
### Blueprint for Service unit test setup
```
func Test{ServiceImpl}_{TestFunc}(t *testing.T) {
    t.Parallel()

	tsi := NewTestService(t)
	defer tsi.Clean()

    type fields struct {
		Ctx     context.Context
		Logger  *zerolog.Logger
		Config  *config.PostServiceConfig
		Service service.Service
	}
    type args struct {
		ctx     context.Context
		...     ...
	}

    type TC struct {
		name      string
		fields    fields
		args      args
		prepare   func(tt *TC)
		validate  func(tt *TC)
		wantErr   bool
		timestamp time.Time
		fixture   interface{}
	}

    tests := []TC{
        {
            name: "Test Sample",
            fields: fields{
				Ctx:     context.TODO(),
				Logger:  &zerolog.Logger{},
				Config:  config.GetConfigFromFile().AppConfig.ServiceConfig.{ServiceConfig},
				Service: tsi.Service,
			},
            args: args{
				ctx: context.TODO(),
                ...: ...
			},
            wantErr: false,
            prepare: func(tt *TC) {},
            validate: func(tt *TC) {},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.timestamp = fakeUTCNow
			pi := &service.{ServiceImpl}{
				Ctx:     tt.fields.Ctx,
				Logger:  tt.fields.Logger,
				Config:  tt.fields.Config,
				Service: tt.fields.Service,
			}
			tt.prepare(&tt)
			got, err := pi.{TestFunc}(tt.args.ctx, ...)
			if (err != nil) != tt.wantErr {
				t.Errorf("{ServiceImpl}.{TestFunc}() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.validate(&tt, got)
        })
    }
}
```

### Clear collection
```
tsi.Service.MongoDB().Cli().Database(model.DBName).Collection(model.CollName).DeleteMany(context.TODO(), bson.M{})
```

### Helper Assert functions
- Assert_TimestampDuration
	```
	allowedDurationDeviation := 1 * time.Second
	Assert_TimestampDuration(t, t1, t2, &allowedDurationDeviation)
	```
- Assert_DocCount
	```
	Assert_DocCount(
					t,
					tt.fields.Service.MongoDB().Cli().Database(model.DBName).Collection(model.CollName),
					bson.M{"_id": resp.ID},
					1,
				)
	```
- Get_DocByFilter
	```
	Get_DocByFilter(tt.fields.Service.MongoDB().Cli().Database(model.DBName).Collection(model.CollName), bson.M{"_id": resp.ID}, &doc)
	```
- Assert_DocJson
	```
	Assert_DocJson(t, doc1, doc2)
	```
