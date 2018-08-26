import consts from "./consts"
import axios from "axios"

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
    console.log('errro!')
    var er = JSON.parse(JSON.stringify(error))
    console.log(er);
    

    console.log(error)      
    console.log(error.response.status) 
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