package app

type ServiceFactory struct {
	xUser  *XUserService
	report *ReportService
}

func (s *ServiceFactory) XUser() *XUserService {
	return s.xUser
}

func (s *ServiceFactory) Report() *ReportService {
	return s.report
}

var instance *ServiceFactory = &ServiceFactory{
	xUser:  NewXUserService(),
	report: NewReportService(),
}

func Factory() *ServiceFactory {
	return instance
}
