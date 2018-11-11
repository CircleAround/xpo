package app

type ServiceFactory struct {
	xUser    *XUserService
	report   *ReportService
	language *LanguageService
}

func (s *ServiceFactory) XUser() *XUserService {
	return s.xUser
}

func (s *ServiceFactory) Report() *ReportService {
	return s.report
}

func (s *ServiceFactory) Language() *LanguageService {
	return s.language
}

var instance *ServiceFactory = &ServiceFactory{
	xUser:    NewXUserService(),
	report:   NewReportService(),
	language: NewLanguageService(),
}

func Factory() *ServiceFactory {
	return instance
}
