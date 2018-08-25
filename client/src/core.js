import consts from "./consts"
import axios from "axios"

const api = axios.create({
  baseURL: consts.API_ENDPOINT,
  headers: {
    "Content-Type": "application/json",
    "X-Requested-With": "XMLHttpRequest"
  },
  responseType: "json"
})

export default {
  status: {
    list: []
  },
  retriveReports() {
    const self = this
    if(self.status.list.length > 0) return Promise.resolve([]);

    return api
      .get("/xreports")
      .then(function(response) {
        response.data.forEach(item => {
          self.status.list.push(item)
        })
      })
      .catch(function(error) {
        console.log(error, "ERROR!! occurred in Backend.")
      })
  }
}