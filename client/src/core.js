import consts from "./consts"
import axios from "axios"
import moment from "moment-timezone"

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
    var er = JSON.parse(JSON.stringify(error))
    if (error.response.status == 401) {
      location.href = "http://localhost:5100"
    }
    else {
      return error
    }
  })
}

export default {
  status: {
    list: []
  },
  retriveReports() {
    if(this.status.list.length > 0) return Promise.resolve([]);

    return errorFilter(api
      .get("/xreports")
      .then((response) => {
        response.data.forEach(item => {
          item.created_at = moment(item.created_at)
          item.updated_at = moment(item.updated_at)
          this.status.list.push(item)
        })
      }))
  },
  postReport(report){
    return errorFilter(api.post('/xreports', report).then((response)=>{
      console.log(response)
      this.status.list.unshift(report)
    }))
  }
}