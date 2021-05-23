package client

import (
	"log"
	"os"
)

// Send interrupt signal and exit smoothly
func (c *Controller) Disconnect() error {
	log.Println("Disconnect from supernode")
	err := c.cmd.Process.Signal(os.Interrupt)
	if err != nil {
		return err
	}
	c.cmd.Wait()
	c.cmd = nil
	return nil
}
