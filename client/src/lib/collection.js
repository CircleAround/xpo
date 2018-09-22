class ListMap {
  constructor(array = []) {
    this.array = array
    this.map = new Map()
    this.updateMap()
  }

  push(object) {
    return this._add(object, 'push')
  }

  unshift(object) {
    return this._add(object, 'unshift')
  }

  enhanceObject(object) {
    // nop
    return object
  }

  getKey(object) {
    throw new Error('Unimplemented method: getKey')
  }

  find(key) {
    return this.map[key]
  }

  updateItem(object) {
    const obj = this.map[this.getKey(object)]
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
    this.array.splice(0, this.array.length)
    this.map.clear()
  }

  _addMap(object) {
    const obj = this.enhanceObject(object)
    this.map[this.getKey(object)] = obj
    return obj
  }

  _add(object, methodName) {
    var obj = this.find(this.getKey(object))
    if (obj) {
      return this._assignItem(obj, object)
    } else {
      obj = this._addMap(object)
      this.array[methodName](obj)
      return obj
    }
  }

  _assignItem(src, dst) {
    Object.assign(dst, this.enhanceObject(src))
    return src
  }
}

export default { ListMap }
