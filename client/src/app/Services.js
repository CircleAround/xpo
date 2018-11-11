function errorFilter(promise) {
  return promise.catch(error => {
    try {
      console.error(error)
      console.error('error', error)

      if (!error.response) {
        throw error
      }

      if (error.response.status === 401) {
        location.href = process.env.API_ENDPOINT
        return
      }

      throw error
    } catch (ex) {
      console.error(error.message)
      console.error('error', error)
      throw ex
    }
  })
}

class UserService {
  constructor(api) {
    this.api = api
  }

  retriveMe() {
    const now = new Date().getTime()
    return this.api.get(`/users/me?_=${now}`)
  }

  getByName(name) {
    return errorFilter(this.api.get(`/users/${name}`))
  }

  getById(id) {
    return errorFilter(this.api.get(`/users/${id}`))
  }

  postXUser(name, nickname) {
    return errorFilter(this.api.post('/users/me', { name, nickname }))
  }

  updateXUser(name, nickname) {
    return errorFilter(this.api.put('/users/me', { name, nickname }))
  }
}

class ReportService {
  constructor(api) {
    this.api = api
  }

  retriveReports() {
    return errorFilter(this.api.get('/reports'))
  }

  searchReportsByAuthorId(authorId) {
    return errorFilter(this.api.get(`/reports/${authorId}`))
  }

  searchReportsByLanguage(language) {
    return errorFilter(this.api.get(`/languages/${language}/reports`))
  }

  searchReportsYmd(authorId, year, month, day) {
    return errorFilter(
      this.api.get(`/reports/${authorId}/_/${year}/${month}/${day}`)
    )
  }

  findReport(authorId, id) {
    return errorFilter(this.api.get(`/reports/${authorId}/${id}`))
  }

  postReport(report) {
    return errorFilter(this.api.post('/reports', report))
  }

  updateReport(report, params) {
    return errorFilter(
      this.api.put(`/reports/${params.authorId}/${params.id}`, report)
    )
  }
}

class LanguageService {
  constructor(api) {
    this.api = api
  }

  getAll() {
    return errorFilter(this.api.get(`/languages`))
  }

  getAllNames() {
    return errorFilter(this.api.get(`/languages/names`))
  }
}

class Services {
  constructor(api) {
    this.users = new UserService(api)
    this.reports = new ReportService(api)
    this.languages = new LanguageService(api)
  }
}

export default Services
