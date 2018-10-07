class ErrorHandler {
  constructor(i18n) {
    this.i18n = i18n
  }

  eachInResponse(error, handler) {
    const i18n = this.i18n
    const e = error.response.data.error

    if (e.type === 'validation') {
      const items = e.items
      Object.keys(items).forEach(property => {
        const item = items[property]
        item.reasons.forEach(reason => {
          this.handleValidation(handler, e.type, reason, property)
        })
      })
      return
    }

    this.handle(handler, e) || this.handleUnexpected(handler, e)
  }

  validation_message(reason, property) {
    return this.i18n.t(this.reason2ValidationMessageKey(reason, property), {
      property: property
    })
  }

  handle(handler, e) {
    const key = this.type2MessageKey(e.type, e.property)
    const msg = this.i18n.t(key, {
      property: e.property
    })
    if (!msg) {
      console.error(`${key} is not found in i18n`)
      return false
    }
    handler(msg, e.type, e.property)
    return true
  }

  handleValidation(handler, type, reason, property) {
    handler(
      this.validation_message(reason, property) ||
        this.validation_message('nothing'),
      type,
      property
    )
  }

  handleUnexpected(handler, e) {
    handler(
      this.i18n.t(`error.messages.unexpected`),
      e.type || 'unexpectedError',
      null
    )
  }

  reason2ValidationMessageKey(reason, property) {
    return `error.messages.validation.${reason}`
  }

  type2MessageKey(type, property) {
    return `error.messages.${e.type}`
  }
}

export default ErrorHandler
