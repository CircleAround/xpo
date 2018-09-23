import axios from 'axios'
import consts from './consts'

const api = axios.create({
  baseURL: consts.API_ENDPOINT,
  headers: {
    'Content-Type': 'application/json',
    'X-Requested-With': 'XMLHttpRequest'
  },
  responseType: 'json',
  withCredentials: true
})

function errorFilter(promise) {
  return promise.catch(error => {
    try {
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
      console.error(ex)
      throw ex
    }
  })
}

const UserService = {
  retriveMe() {
    return api.get('/users/me')
  },
  postXUser(name, nickname) {
    return errorFilter(api.post('/users/me', { name, nickname }))
  }
}

const ReportService = {
  retriveReports() {
    return errorFilter(api.get('/reports'))
  },
  findReport(authorId, id) {
    return errorFilter(api.get(`/reports/${authorId}/${id}`))
  },
  postReport(report) {
    return errorFilter(api.post('/reports', report))
  },
  updateReport(report, params) {
    return errorFilter(
      api.put(`/reports/${params.user_id}/${params.id}`, report)
    )
  }
}

export default {
  users: UserService,
  reports: ReportService
}