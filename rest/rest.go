package rest

import (
	"encoding/json"
	"fmt"
	"os"

	bolt "github.com/coreos/bbolt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//GRPCCfg aaaa
type GRPCCfg struct {
	Host        string  `json:"host"`
	Port        int     `json:"port"`
	User        string  `json:"user"`
	Password    string  `json:"password"`
	Meta        bool    `json:"meta"`
	EOS         bool    `json:"eos"`
	CID         string  `json:"cid"`
	WS          int32   `json:"ws"`
	TLS         TLSCfg  `json:"tls"`
	Paths       []Spath `json:"paths"`
	Compression string  `json:"compression"`
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

//Cfg is an HTTP cpnfig struct
type Cfg struct {
	Addr, Port, Secret string
}

type handler struct {
	db *bolt.DB
}

func readCfg() Cfg {
	var cfg Cfg
	cfg.Port = os.Getenv("PORT")
	cfg.Addr = os.Getenv("ADDRESS")
	if cfg.Port == "" {
		cfg.Port = "8888"
	}
	return cfg
}

//StartHTTPSrv strts http server
func StartHTTPSrv(db *bolt.DB) error {
	cfg := readCfg()
	h := handler{
		db: db,
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
		api.PUT("/device", h.addDevice)
		api.DELETE("/device/:id", h.delDevice)
	}
	logrus.WithFields(logrus.Fields{
		"Port": cfg.Port,
		"Addr": cfg.Addr,
	}).Info("http server starting ...")
	err := router.Run(cfg.Addr + ":" + cfg.Port)
	if err != nil {
		return err
	}
	return nil
}

func (h *handler) delDevice(c *gin.Context) {
	hostname := c.Param("id")
	fmt.Println(hostname)
	err := h.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("devices"))
		if b != nil {
			err := b.Delete([]byte(hostname))
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
	c.JSON(200, hostname)
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
				fmt.Println(string(k), string(v))
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

func (h *handler) addDevice(c *gin.Context) {
	var d GRPCCfg
	err := c.BindJSON(&d)
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}
	err = h.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("devices"))
		if err != nil {
			return err
		}
		data, err := json.Marshal(&d)
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}
		return b.Put([]byte(d.Host), data)
	})
	if err != nil {
		c.AbortWithStatusJSON(500, err)
	}
	c.JSON(200, &d)
}
