package oauth

// Static clients for PoC. In real systems, load from config.
var staticClients = []Client{
	{
		ID: "mcp-cli-12345",
		AllowedRedirects: []string{
			// Allow 127.0.0.1 or ::1 with any port and fixed path /callback (validated programmatically)
			// For simplicity, we also allow custom scheme example:
			"myapp://oauth2redirect",
		},
		Public: true,
	},
}

func FindClient(id string) *Client {
	for i := range staticClients {
		if staticClients[i].ID == id {
			return &staticClients[i]
		}
	}
	return nil
}
