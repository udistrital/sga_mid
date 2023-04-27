package controllers

import (
	"time"

	"github.com/astaxie/beego"
)

// Time_bogController operations for Time_bog
type Time_bogController struct {
	beego.Controller
}

// URLMapping ...
func (c *Time_bogController) URLMapping() {
	c.Mapping("GetTimeBog", c.GetTimeBog)
}

// GetTimeBog ...
// @Title GetTimeBog
// @Description get Time_bog
// @Success 200 {object} models.Time_bog
// @Failure 500 something bad happened
// @router / [get]
func (c *Time_bogController) GetTimeBog() {
	what_time_is_it := time.Now()
	inUTC, _ := time.LoadLocation("UTC")
	inBog, _ := time.LoadLocation("America/Bogota")
	data := map[string]interface{}{
		"UNIX": what_time_is_it.Unix() * 1000,                  // Unix timestamp fixed to seconds represeted in milliseconds
		"UTC":  what_time_is_it.In(inUTC).Format(time.RFC3339), // UTC timestamp fixed to seconds
		"BOG":  what_time_is_it.In(inBog).Format(time.RFC3339), // BOG timestamp fixed to seconds
	}
	c.Ctx.Output.SetStatus(200)
	c.Data["json"] = map[string]interface{}{"Success": true, "Status": "200", "Message": "Query successful", "Data": data}
	c.ServeJSON()
}
