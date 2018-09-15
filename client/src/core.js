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

function enhanceReport(item) {
  item.created_at = moment(item.created_at)
  item.updated_at = moment(item.updated_at)
  item.markdown = function() {
    return marked(this.content)
  }
  return item
}

export default {
  state: {
    me: {
      id: null,
      email: null,
      name: null,
      login_url: null,
      logout_url: null
    },
    list: [],
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
    return api
      .get('/users/me')
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
    return errorFilter(
      api.post('/users/me', { name, nickname }).then(response => {
        this.state.me = response.data
      })
    )
  },
  retriveReports() {
    return errorFilter(
      api.get('/reports').then(response => {
        response.data.forEach(item => {
          this.state.list.push(enhanceReport(item))
        })
      })
    )
  },
  postReport() {
    return errorFilter(
      api.post('/reports', this.state.newReport).then(response => {
        this.state.list.unshift(enhanceReport(response.data))
        this.posted = true
        this.initNewReport()
        router.push('/')
      })
    )
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
