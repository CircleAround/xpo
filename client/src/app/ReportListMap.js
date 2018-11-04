import moment from 'moment-timezone'
import { ListMap } from '../lib/collection'
import marked from 'marked'
import hljs from 'highlightjs'
import jstimezonedetect from 'jstimezonedetect'

marked.setOptions({
  highlight: function(code, lang) {
    return hljs.highlightAuto(code, [lang]).value
  }
})

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

class ReportListMap extends ListMap {
  getKey(object) {
    return `${object.authorId}/${object.id}`
  }

  enhanceObject(object) {
    return enhanceReport(object)
  }
}

export default ReportListMap
