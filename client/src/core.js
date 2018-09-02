import consts from './consts'
import axios from 'axios'
import moment from 'moment-timezone'
import marked from 'marked'
import router from './router'

moment.tz.setDefault('Asia/Tokyo')

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
  return promise.catch((error) => {
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

function enhanceReport(item) {
  item.created_at = moment(item.created_at)
  item.updated_at = moment(item.updated_at)
  item.markdown = function () { return marked(this.content) }
  return item
}

export default {
  status: {
    list: [],
    newReport: {},
    posted: false
  },
  initialize() {
    this.retriveReports().catch(function (error) {
      console.log(error)
    })

    this.initNewReport()
  },
  retriveReports() {
    return errorFilter(api
      .get('/reports')
      .then((response) => {
        response.data.forEach(item => {
          this.status.list.push(enhanceReport(item))
        })
      }))
  },
  postReport() {
    return errorFilter(api.post('/reports', this.newReport).then((response) => {
      this.status.list.unshift(enhanceReport(response.data))
      this.posted = true
      this.initNewReport()
      router.push('/')
    }))
  },
  initNewReport() {
    this.newReport = {
      content: '',
      'content-type': 'text/x-markdown'
    }
  },
  getMessagesOfValidationError(error) {
    if (error.response.data.error.type === 'ValidationError') {
      const items = error.response.data.error.items
      let ret = []
      Object.keys(items).forEach(property => {
        const item = items[property]
        ret = item.reasons
          .map(reason => {
            if (reason === 'required') {
              return `${property}は必須です`
            }
            if (reason === 'toolong') {
              return `${property}は長すぎます`
            }
            return `${property}が何らかのエラーです`
          })
      })
      return ret
    }
  }
}
