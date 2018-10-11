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
    return this.api.get('/users/me')
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

class Services {
  constructor(api) {
    this.users = new UserService(api)
    this.reports = new ReportService(api)
  }
}

export default Services
