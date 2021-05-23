package client

// On Windows can only kill
func (c *Controller) Disconnect() error {
	err := c.cmd.Process.Kill()
	if err != nil {
		return err
	}
	c.cmd.Wait()
	c.cmd = nil
	return nil
}
