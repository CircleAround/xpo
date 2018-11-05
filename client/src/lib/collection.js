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

  pop() {
    return this._removeMap(this.array.pop())
  }

  shift() {
    return this._removeMap(this.array.shift())
  }

  unshift(object) {
    return this._add(object, 'unshift')
  }

  replaceAt(index, object) {
    const ret = this._removeMap(this.array[index])
    this.array.splice(index, 1, object)
    this._addMap(object)
    return ret
  }

  // forceUpdate for send update event on vue.js
  forceUpdate() {
    this.array.push(this.array.pop())
  }

  enhanceObject(object) {
    // nop for override
    return object
  }

  getKey(object) {
    throw new Error('Unimplemented method: getKey')
  }

  at(index) {
    return this.array[index]
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

  reset(objects) {
    this.clear()
    this.pushAll(objects)
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

  _removeMap(object) {
    this.map.delete(this.getKey(object))
    return object
  }
}

function remove(array, item) {
  array.splice(array.indexOf(item), 1)
}

export { ListMap, remove }
