package app

type ServiceFactory struct {
	xUser    *XUserService
	report   *ReportService
	language *LanguageService
	project  *ProjectService
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

func (s *ServiceFactory) Project() *ProjectService {
	return s.project
}

var instance *ServiceFactory = &ServiceFactory{
	xUser:    NewXUserService(),
	report:   NewReportService(),
	language: NewLanguageService(),
	project:  NewProjectService(),
}

func Factory() *ServiceFactory {
	return instance
}
