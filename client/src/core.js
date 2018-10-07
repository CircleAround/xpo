import moment from 'moment-timezone'
import marked from 'marked'
import jstimezonedetect from 'jstimezonedetect'
import router from './router'
import service from './service'
import collection from './lib/collection'

var tz = jstimezonedetect.determine()
moment.tz.setDefault(tz.name())

function enhanceReport(item) {
  if (item.markdown) {
    // already enhanced guard
    return item
  }

  item.reportedAt = moment(item.reportedAt)
  item.createdAt = moment(item.createdAt)
  item.updatedAt = moment(item.updatedAt)
  item.markdown = function() {
    return marked(this.content)
  }
  return item
}

class ReportListMap extends collection.ListMap {
  getKey(object) {
    return `${object.authorId}/${object.id}`
  }

  enhanceObject(object) {
    return enhanceReport(object)
  }
}

const listMap = new ReportListMap()
const subListMap = new ReportListMap()

export default {
  state: {
    me: {
      id: null,
      email: null,
      name: null,
      loginUrl: null,
      logoutUrl: null
    },
    list: listMap.array,
    subList: subListMap.array,
    newReport: { content: null },
    posted: false
  },
  initialize() {
    this.initNewReport()
    return Promise.all([this.retriveMe(), this.retriveReports()])
  },
  isLoggedIn() {
    if (!this.state) return false
    return this.state.me.id != null
  },
  retriveMe: async function() {
    try {
      const response = await service.users.retriveMe()
      if (response.data === 'BE_SIGN_UP') {
        router.push('/signup')
      } else {
        this.state.me = response.data
      }
    } catch (error) {
      if (!error.response) {
        throw error
      }

      if (error.response.status !== 401) {
        throw error
      }

      this.state.me.loginUrl = error.response.data.error
    }
  },
  postXUser: async function(name, nickname) {
    const response = await service.users.postXUser(name, nickname)
    this.state.me = response.data
  },
  updateXUser: async function(name, nickname) {
    const response = await service.users.updateXUser(name, nickname)
    this.state.me = response.data
  },
  retriveReports: async function() {
    listMap.clear()
    const response = await service.reports.retriveReports()
    listMap.pushAll(response.data)
  },
  searchReportsYmd: async function(authorId, year, month, day) {
    subListMap.clear()
    const response = await service.reports.searchReportsYmd(
      authorId,
      year,
      month,
      day
    )
    subListMap.pushAll(response.data)
  },
  findReport: async function(authorId, id) {
    const response = await service.reports.findReport(authorId, id)
    return listMap.push(response.data)
  },
  findReport4Update: async function(authorId, id) {
    const newObject = await this.findReport(authorId, id)
    Object.assign(this.state.newReport, newObject)
  },
  postReport: async function() {
    const response = await service.reports.postReport(this.state.newReport)
    listMap.unshift(response.data)
    this.initNewReport()
    router.push('/')
  },
  updateReport: async function(params) {
    const response = await service.reports.updateReport(
      this.state.newReport,
      params
    )
    listMap.updateItem(response.data)
    this.initNewReport()
    router.push('/')
  },
  initNewReport() {
    this.state.newReport = {
      content: '',
      contentType: 'text/x-markdown'
    }
  },
  forceUpdateMainList() {
    listMap.forceUpdate()
  }
}
