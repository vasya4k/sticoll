package rest

import (
	"encoding/json"
	"sync"

	bolt "github.com/coreos/bbolt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const zeroUUID = "00000000-0000-0000-0000-000000000000"

//GRPCCfg aaaa
type GRPCCfg struct {
	Host        string    `json:"host"`
	Port        int       `json:"port"`
	User        string    `json:"user"`
	Password    string    `json:"password"`
	Meta        bool      `json:"meta"`
	EOS         bool      `json:"eos"`
	CID         string    `json:"cid"`
	WS          int32     `json:"ws"`
	TLS         TLSCfg    `json:"tls"`
	Paths       []Spath   `json:"paths"`
	Compression string    `json:"compression"`
	UUID        uuid.UUID `json:"uuid"`
	Removed     bool      `json:"removed"`
	sync.RWMutex
}

//Spath aaa
type Spath struct {
	Path string `json:"path"`
	Freq uint64 `json:"freq"`
	Mode string `json:"mode"`
}

//TLSCfg aaa
type TLSCfg struct {
	Enabled    bool   `json:"enabled"`
	ClientCrt  string `json:"client_crt"`
	ClientKey  string `json:"client_key"`
	CA         string `json:"ca"`
	ServerName string `json:"server_name"`
}

//HTTPCfg holds http configuration gets passtes into this pkg
type HTTPCfg struct {
	Port   string
	UIPath string
	Addr   string
}

type handler struct {
	db    *bolt.DB
	cfgs  *[]*GRPCCfg
	cfgCh chan *GRPCCfg
}

//StartHTTPSrv strts http server
func StartHTTPSrv(hcfg *HTTPCfg, db *bolt.DB, cfgs *[]*GRPCCfg, cfgCh chan *GRPCCfg) error {
	h := handler{
		db:    db,
		cfgs:  cfgs,
		cfgCh: cfgCh,
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowHeaders = []string{"*"}
	// config.AllowOrigins = []string{"*"}
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"*"}
	config.ExposeHeaders = []string{"*"}
	//Need cors if UI is served from a different server
	router.Use(cors.New(config))
	//Add a version 1 group
	api := router.Group("/v1")

	{
		api.GET("/devices", h.getDevices)
		api.POST("/device", h.addDevice)
		api.PUT("/device", h.updDevice)
		api.DELETE("/device/:id", h.delDevice)
	}
	logrus.WithFields(logrus.Fields{
		"Port": hcfg.Port,
		"Addr": hcfg.Addr,
	}).Info("http server starting ...")
	err := router.Run(hcfg.Addr + ":" + hcfg.Port)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) delDevice(c *gin.Context) {
	ud, err := uuid.FromString(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}
	err = h.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))
		if b != nil {
			err := b.Delete(ud.Bytes())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}
	for _, cfg := range *h.cfgs {
		cfg.Lock()
		if cfg.UUID == ud {
			cfg.Removed = true
		}
		cfg.Unlock()
	}
	c.JSON(200, c.Param("id"))
}

func (h *handler) getDevices(c *gin.Context) {
	gCfgs := make([]*GRPCCfg, 0)
	err := h.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))
		if b != nil {
			err := b.ForEach(func(k, v []byte) error {
				var g GRPCCfg
				err := json.Unmarshal(v, &g)
				if err != nil {
					return err
				}
				gCfgs = append(gCfgs, &g)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.AbortWithStatusJSON(500, err.Error())
		return
	}
	c.JSON(200, gCfgs)
}

func (h *handler) updDevice(c *gin.Context) {
	var d GRPCCfg
	err := c.BindJSON(&d)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	err = h.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("devices"))
		if err != nil {
			return err
		}
		data, err := json.Marshal(&d)
		if err != nil {
			return err
		}
		return b.Put(d.UUID.Bytes(), data)
	})
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	c.JSON(200, &d)
}

func (h *handler) addDevice(c *gin.Context) {
	var d GRPCCfg
	err := c.BindJSON(&d)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	d.UUID = uuid.NewV4()
	err = h.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("devices"))
		if err != nil {
			return err
		}
		data, err := json.Marshal(&d)
		if err != nil {
			return err
		}
		return b.Put(d.UUID.Bytes(), data)
	})
	if err != nil {
		c.AbortWithStatusJSON(500, err)
		return
	}
	h.cfgCh <- &d
	c.JSON(200, &d)
}
