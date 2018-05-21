package pm2_io_apm_go

import "testing"

func Connect(t *testing.T) {
	Pm2Io := Pm2Io{
		PublicKey:  "9nc25845w31vqeq",
		PrivateKey: "1e34mwmtaid0pr7",
	}

	t.Run("Start", func(t *testing.T) {
		Pm2Io.Start()
		Pm2Io.SendStatus()
	})
}
