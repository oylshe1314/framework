package server

type Manager interface {
	Init() error
	Close() error
}

func InitManagers(mgrs ...Manager) (err error) {
	for _, mgr := range mgrs {
		err = mgr.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func CloseManagers(mgrs ...Manager) (errs []error) {
	errs = make([]error, len(mgrs))
	for i, mgr := range mgrs {
		errs[i] = mgr.Close()
	}
	return errs
}
