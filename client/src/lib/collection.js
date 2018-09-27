class ListMap {
  constructor(array = []) {
    this.array = array
    this.map = new Map()
    this.updateMap()
  }

  push(object) {
    return this._add(object, 'push')
  }

  pushAll(objects) {
    objects.forEach(object => {
      this.push(object)
    })
  }

  unshift(object) {
    return this._add(object, 'unshift')
  }

  enhanceObject(object) {
    // nop for override
    return object
  }

  getKey(object) {
    throw new Error('Unimplemented method: getKey')
  }

  find(key) {
    return this.map.get(key)
  }

  updateItem(object) {
    const obj = this.map.get(this.getKey(object))
    if (!obj) return
    this._assignItem(object, obj)
  }

  updateMap() {
    this.map.clear()
    this.array.forEach(object => {
      this._addMap(object)
    })
  }

  clear() {
    this.array.splice(0)
    this.map.clear()
  }

  _addMap(object) {
    const obj = this.enhanceObject(object)
    this.map.set(this.getKey(object), obj)
    return obj
  }

  _add(object, methodName) {
    var obj = this.find(this.getKey(object))
    if (obj) {
      return this._assignItem(object, obj)
    } else {
      obj = this._addMap(object)
      this.array[methodName](obj)
      return obj
    }
  }

  _assignItem(src, dst) {
    Object.assign(dst, this.enhanceObject(src))
    return dst
  }
}

export default { ListMap }
