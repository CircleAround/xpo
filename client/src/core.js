import moment from 'moment-timezone'
import marked from 'marked'
import router from './router'
import service from './service'
import collection from './lib/collection'

moment.tz.setDefault('Asia/Tokyo')

function enhanceReport(item) {
  item.created_at = moment(item.created_at)
  item.updated_at = moment(item.updated_at)
  item.markdown = function() {
    return marked(this.content)
  }
  return item
}

class ReportListMap extends collection.ListMap {
  getKey(object) {
    return `${object.author_id}/${object.id}`
  }

  enhanceObject(object) {
    return enhanceReport(object)
  }
}

const listMap = new ReportListMap()

export default {
  state: {
    me: {
      id: null,
      email: null,
      name: null,
      login_url: null,
      logout_url: null
    },
    list: listMap.array,
    newReport: { content: null },
    posted: false
  },
  initialize() {
    this.retriveMe()
    this.retriveReports().catch(function(error) {
      console.log(error)
    })

    this.initNewReport()
  },
  isLoggedIn() {
    if (!this.state) return false
    return this.state.me.id != null
  },
  retriveMe() {
    return service.users
      .retriveMe()
      .then(response => {
        if (response.data === 'BE_SIGN_UP') {
          router.push('/signup')
        } else {
          this.state.me = response.data
        }
      })
      .catch(error => {
        if (!error.response) {
          throw error
        }

        if (error.response.status !== 401) {
          throw error
        }

        this.state.me.login_url = error.response.data.error
      })
  },
  postXUser(name, nickname) {
    return service.users.postXUser(name, nickname).then(response => {
      this.state.me = response.data
    })
  },
  retriveReports() {
    return service.reports.retriveReports().then(response => {
      response.data.forEach(item => {
        listMap.push(item)
      })
    })
  },
  findReport(authorId, id) {
    return service.reports.findReport(authorId, id).then(response => {
      return listMap.push(response.data)
    })
  },
  findReport4Update(authorId, id) {
    return this.findReport(authorId, id).then(newObject => {
      Object.assign(this.state.newReport, newObject)
    })
  },
  postReport() {
    return service.reports.postReport(this.state.newReport).then(response => {
      listMap.unshift(response.data)
      this.initNewReport()
      router.push('/')
    })
  },
  updateReport(params) {
    return service.reports
      .updateReport(this.state.newReport, params)
      .then(response => {
        listMap.updateItem(response.data)
        this.initNewReport()
        router.push('/')
      })
  },
  initNewReport() {
    this.state.newReport = {
      content: '',
      'content-type': 'text/x-markdown'
    }
  },
  eachResponseErrors(error, handler) {
    const e = error.response.data.error

    const i18n = {
      required: property => {
        return `${property}は必須です`
      },
      toolong: property => {
        return `${property}は長すぎます`
      },
      username_format: property => {
        return `${property}に利用できる文字は半角英数小文字です`
      },
      usernickname_format: property => {
        return `ニックネームに利用できる文字に一致しませんでした`
      },
      nothing: property => {
        return `${property}が何らかのエラーです`
      }
    }

    switch (e.type) {
      case 'ValidationError':
        const items = e.items
        Object.keys(items).forEach(property => {
          const item = items[property]
          item.reasons.forEach(reason => {
            if (i18n[reason]) {
              return handler(i18n[reason](property), e.type, property)
            }
            handler(i18n['nothing'](property), e.type, property)
          })
        })
        break
      case 'ValueNotUniqueError':
        handler(`${e.property}は既に存在します`, e.type, e.property)
        break
      case 'DuplicatedObjectError':
        handler(`既に存在します`, e.type, null)
        break

      default:
        break
    }
  }
}
