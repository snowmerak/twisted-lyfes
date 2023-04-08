package redis

func (c *Client) AddPeerToSet(name, peer string) error {
	return c.client.Do(c.ctx, c.client.B().Set().Key(name).Value(peer).Build()).Error()
}
