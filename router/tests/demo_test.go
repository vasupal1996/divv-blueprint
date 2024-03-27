package router_test

import (
	"bytes"
	"go-app/internals/config"
	"go-app/router"
	"go-app/schema"
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestRouter_HelloWorldHandler(t *testing.T) {

	tri := NewRouterTest(t)
	defer tri.Clean()

	type fields struct {
		App    *fiber.App
		Logger *zerolog.Logger
		Config *config.RouterConfig
	}

	type args struct {
		c *fiber.Ctx
	}

	type TC struct {
		name          string
		url           string
		method        string
		body          io.Reader
		fields        fields
		args          args
		prepare       func(tt *TC)
		checkResponse func(tt *TC, resp *http.Response)
	}

	tests := []TC{
		{
			name:   "Test Success",
			url:    "/",
			method: http.MethodGet,
			body:   nil,
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {
				tri.demoService.EXPECT().DemoFunc(gomock.Any()).Return("test").Times(1)
			},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, `{"message":"Hello World"}`, string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router.Router{
				App:         tt.fields.App,
				Logger:      tt.fields.Logger,
				Config:      tt.fields.Config,
				DemoService: tri.demoService,
			}
			r.RegisterRoutes()
			tt.prepare(&tt)
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			resp, err := r.App.Test(req)
			assert.Nil(t, err)
			tt.checkResponse(&tt, resp)
		})
	}
}

func TestRouter_BadRequestHandler(t *testing.T) {

	tri := NewRouterTest(t)
	defer tri.Clean()

	type fields struct {
		App    *fiber.App
		Logger *zerolog.Logger
		Config *config.RouterConfig
	}

	type args struct {
		c *fiber.Ctx
	}

	type TC struct {
		name          string
		url           string
		method        string
		body          io.Reader
		fields        fields
		args          args
		prepare       func(tt *TC)
		checkResponse func(tt *TC, resp *http.Response)
	}

	tests := []TC{
		{
			name:   "Test Success",
			url:    "/bad-request",
			method: http.MethodGet,
			body:   nil,
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, `{"success":false,"error":[{"code":"BadRequest","msg":"request failed"}]}`, string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router.Router{
				App:         tt.fields.App,
				Logger:      tt.fields.Logger,
				Config:      tt.fields.Config,
				DemoService: tri.demoService,
			}
			r.RegisterRoutes()
			tt.prepare(&tt)
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			resp, err := r.App.Test(req)
			assert.Nil(t, err)
			tt.checkResponse(&tt, resp)
		})
	}
}

func TestRouter_InternalServerHandler(t *testing.T) {

	tri := NewRouterTest(t)
	defer tri.Clean()

	type fields struct {
		App    *fiber.App
		Logger *zerolog.Logger
		Config *config.RouterConfig
	}

	type args struct {
		c *fiber.Ctx
	}

	type TC struct {
		name          string
		url           string
		method        string
		body          io.Reader
		fields        fields
		args          args
		prepare       func(tt *TC)
		checkResponse func(tt *TC, resp *http.Response)
	}

	tests := []TC{
		{
			name:   "Test Success",
			url:    "/internal-error",
			method: http.MethodGet,
			body:   nil,
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.Equal(t, `runtime error: index out of range [0] with length 0`, string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router.Router{
				App:         tt.fields.App,
				Logger:      tt.fields.Logger,
				Config:      tt.fields.Config,
				DemoService: tri.demoService,
			}
			r.App.Use(recover.New(recover.Config{
				EnableStackTrace: false,
			}))
			r.RegisterRoutes()
			tt.prepare(&tt)
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			resp, err := r.App.Test(req)
			assert.Nil(t, err)
			tt.checkResponse(&tt, resp)
		})
	}
}

func TestRouter_InsertOneHandler(t *testing.T) {

	tri := NewRouterTest(t)
	defer tri.Clean()

	type fields struct {
		App    *fiber.App
		Logger *zerolog.Logger
		Config *config.RouterConfig
	}

	type args struct {
		c *fiber.Ctx
	}

	type TC struct {
		name          string
		url           string
		method        string
		body          io.Reader
		fields        fields
		args          args
		prepare       func(tt *TC)
		checkResponse func(tt *TC, resp *http.Response)
	}

	tests := []TC{
		{
			name:   "No Request Body",
			url:    "/insert",
			method: http.MethodPost,
			body:   nil,
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, `{"success":false,"error":[{"code":"StatusBadRequest","msg":"Request body must not be empty"}]}`, string(data))
			},
		},
		{
			name:   "Invalid Request Body",
			url:    "/insert",
			method: http.MethodPost,
			body:   bytes.NewReader([]byte("invalid body string")),
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, `{"success":false,"error":[{"code":"StatusBadRequest","msg":"Request body contains badly-formed JSON (at position 0)"}]}`, string(data))
			},
		},
		{
			name:   "Success",
			url:    "/insert",
			method: http.MethodPost,
			body:   bytes.NewBuffer([]byte(`{"name":"Lorem"}`)),
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {
				oid, _ := primitive.ObjectIDFromHex("6602ef6e0dc2f69705594eb3")
				tri.demoService.EXPECT().
					InsertOne(gomock.Any(), &schema.InsertOneOpts{Name: "Lorem"}).
					Return(oid, nil).
					Times(1)
			},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusOK, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, `{"success":true,"payload":{"id":"6602ef6e0dc2f69705594eb3"}}`, string(data))
			},
		},
		{
			name:   "Validation Error",
			url:    "/insert",
			method: http.MethodPost,
			body:   bytes.NewBuffer([]byte(`{"name":""}`)),
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {

			},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, `{"success":false,"error":[{"code":"ValidationErr","msg":"name is a required field","field":"name"}]}`, string(data))
			},
		},
		{
			name:   "Validation Error Multiple",
			url:    "/insert",
			method: http.MethodPost,
			body:   bytes.NewBuffer([]byte(`{"name":"LoremIpsumLoremIpsum"}`)),
			fields: fields{
				App:    tri.App,
				Logger: tri.Logger,
				Config: tri.Config,
			},
			args: args{
				c: &fiber.Ctx{},
			},
			prepare: func(tt *TC) {

			},
			checkResponse: func(tt *TC, resp *http.Response) {
				assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
				data, err := io.ReadAll(resp.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, `{"success":false,"error":[{"code":"ValidationErr","msg":"name should not be equal to LoremIpsumLoremIpsum","field":"name"}]}`, string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &router.Router{
				App:         tt.fields.App,
				Logger:      tt.fields.Logger,
				Config:      tt.fields.Config,
				DemoService: tri.demoService,
				Validator:   router.NewValidator(),
			}
			r.RegisterRoutes()
			tt.prepare(&tt)
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			req.Header.Add("Content-Type", "application/json")
			assert.Nil(t, err)
			resp, err := r.App.Test(req)
			assert.Nil(t, err)
			tt.checkResponse(&tt, resp)
		})
	}
}
