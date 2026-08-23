import '@testing-library/jest-dom/vitest'

// jsdom abhi <dialog> ke showModal/close implement nahi karta — polyfill.
if (typeof HTMLDialogElement.prototype.showModal !== 'function') {
  HTMLDialogElement.prototype.showModal = function () {
    this.setAttribute('open', '')
  }
}
if (typeof HTMLDialogElement.prototype.close !== 'function') {
  HTMLDialogElement.prototype.close = function () {
    this.removeAttribute('open')
    this.dispatchEvent(new Event('close'))
  }
}
