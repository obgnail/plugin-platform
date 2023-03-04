package handler

func InitHostBoot() error {
	h := Default()
	h.Run()
	return nil
}
