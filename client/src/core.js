import consts from "./consts"
import axios from "axios"
import moment from "moment-timezone"
import events from 'events'
import marked from 'marked'
import router from './router'

moment.tz.setDefault("Asia/Tokyo");

const api = axios.create({
  baseURL: consts.API_ENDPOINT,
  headers: {
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest"
  },
  responseType: "json",
  withCredentials: true
})

function errorFilter(promise){
  return promise.catch((error)=>{
    eventEmitter.emit('error', error)
    console.error('error', error)

    if(!error.response) {
      return error
    }

    if (error.response.status == 401) {
      location.href = process.env.API_ENDPOINT
      return
    }
    
    return error
  })
}

function enhanceReport(item) {
  item.created_at = moment(item.created_at)
  item.updated_at = moment(item.updated_at)
  item.markdown = function(){ return marked(this.content) }
  return item
}

const EventEmitter = events.EventEmitter;
const eventEmitter = new EventEmitter();

export default {
  status: {
    list: [],
    posted: false,
    eventEmitter: eventEmitter
  },
  retriveReports() {
    if(this.status.list.length > 0) return Promise.resolve([]);

    return errorFilter(api
      .get("/xreports")
      .then((response) => {
        response.data.forEach(item => {
          this.status.list.push(enhanceReport(item))
        })
      }))
  },
  postReport(report){
    return errorFilter(api.post('/xreports', report).then((response)=>{
      this.status.list.unshift(enhanceReport(report))
      this.posted = true
      router.push('/')
    }))
  }
}