import consts from '../consts'
import axios from 'axios'
import Services from './Services'

class DefaultServiceFactory {
  create() {
    const api = axios.create({
      baseURL: consts.API_ENDPOINT,
      headers: {
        'Content-Type': 'application/json',
        'X-Requested-With': 'XMLHttpRequest'
      },
      responseType: 'json',
      withCredentials: true
    })
    return new Services(api)
  }
}

export default DefaultServiceFactory
