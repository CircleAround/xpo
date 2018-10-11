class ErrorHandler {
  constructor(i18n) {
    this.i18n = i18n
  }

  eachInResponse(data, handler) {
    const e = data.error

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

  validationMessage(reason, property) {
    return this.i18n.t(this.validationMessageKeyByReason(reason, property), {
      property: property
    })
  }

  handle(handler, e) {
    const key = this.messageKeyByType(e.type, e.property)
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
      this.validationMessage(reason, property) ||
        this.validationMessage('nothing'),
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

  validationMessageKeyByReason(reason, property) {
    return `error.messages.validation.${reason}`
  }

  messageKeyByType(type, property) {
    return `error.messages.${type}`
  }
}

export default ErrorHandler
