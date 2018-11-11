import router from './router'
import ReportListMap from './app/ReportListMap'
import 'highlightjs/styles/agate.css'
import { remove } from './lib/collection'

const listMap = new ReportListMap()
const subListMap = new ReportListMap()

var service
export function setServices(s) {
  service = s
}

const core = {
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
    targetReport: { content: null },
    languages: [],
    languageNames: [],
    alerts: [],
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
  getXUserByName: async function(name) {
    return service.users.getByName(name)
  },
  getXUserById: async function(id) {
    return service.users.getById(id)
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
  searchReportsByAuthor: async function(authorId) {
    subListMap.clear()
    const response = await service.reports.searchReportsByAuthorId(authorId)
    subListMap.pushAll(response.data)
  },
  searchReportsByLanguage: async function(language) {
    subListMap.clear()
    const response = await service.reports.searchReportsByLanguage(language)
    subListMap.pushAll(response.data)
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
    Object.assign(this.state.targetReport, newObject)
  },
  postReport: async function() {
    const response = await service.reports.postReport(this.state.targetReport)
    listMap.unshift(response.data)
    this.initNewReport()
    router.push('/')
  },
  updateReport: async function(params) {
    const response = await service.reports.updateReport(
      this.state.targetReport,
      params
    )
    listMap.updateItem(response.data)
    this.initNewReport()
    router.push('/')
  },
  initNewReport() {
    this.state.targetReport = {
      content: '',
      contentType: 'text/x-markdown'
    }
  },
  removeLanguageOnTargetReport(name) {
    remove(this.state.targetReport.languages, name)
  },
  getAllLanguageNames: async function() {
    const response = await service.languages.getAllNames()
    response.data.forEach(lng => {
      this.state.languageNames.push(lng)
    })
  },
  getAllLanguages: async function() {
    if (this.state.languages.length) return

    const response = await service.languages.getAll()
    response.data.forEach(lng => {
      this.state.languages.push(lng)
    })
  },
  forceUpdateMainList() {
    listMap.forceUpdate()
  },
  alert(message, type) {
    const msg = { message: message, type: type || 'info' }
    this.state.alerts.push(msg)
    setTimeout(() => {
      remove(this.state.alerts, msg)
    }, 2000)
  }
}

export default core
