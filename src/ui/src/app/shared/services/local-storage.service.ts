import { Injectable } from '@angular/core';
import { isNil } from 'lodash';

const KEY_PREFIX = 'fc';

@Injectable()
export class LocalStorageService {

  constructor() { }

  /**
   * Returns native store provider
   * @returns {Storage}
   */
  get store() {
    return window.localStorage;
  }

  /**
   * Get prefixed key name
   * @param {string} key
   * @returns {string}
   */
  getKeyNamePrefixed(key: string) {
    return `${KEY_PREFIX}.${key}`;
  }

  /**
   * Checks if the key exists in the store
   * @param {string} key
   * @returns {boolean}
   */
  hasKey(key: string): boolean {
    const prefName = this.getKeyNamePrefixed(key);
    return this.store.getItem(prefName) !== null;
  }

  /**
   * Get value from the storage
   * @param {string} key
   * @param otherwise
   * @param {boolean} unprefixed Don't add default prefix to the key name
   * @returns {string}
   */
  getItem(key: string, otherwise: any = null, unprefixed = false): string {
    const prefName = unprefixed ? key : this.getKeyNamePrefixed(key);
    const val = this.store.getItem(prefName);

    return isNil(val) ? otherwise : val;
  }

  /**
   * Get array or object from the store
   * @param {string} key Key
   * @param otherwise Default value
   * @param {boolean} unprefixed Don't add default prefix to the key name
   * @returns {any}
   */
  getObject(key: string, otherwise: any = null, unprefixed = false): any {
    const raw = this.getItem(key, otherwise, unprefixed);

    try {
      return JSON.parse(raw);
    } catch (ex) {
      return otherwise;
    }

  }

  /**
   * Add data to the storage (value will be serialized automatically)
   * @param {string} key
   * @param data Data
   * @param unprefixed Don't add default prefix
   */
  setItem(key: string, data: any, unprefixed = false) {
    let serialized = null;
    const pref = unprefixed ? key : this.getKeyNamePrefixed(key);

    switch (typeof data) {
      case 'number':
        serialized = String(data);
        break;

      case 'object':
        serialized = JSON.stringify(data);
        break;

      case 'function':
        throw new TypeError('Function values cannot be serialized');

      default:
        serialized = data;
        break;
    }

    this.store.setItem(pref, serialized);
  }

  /**
   * Remove item from the store
   * @param {string} key Key
   * @param {boolean} unprefixed Don't add default key prefix
   */
  removeItem(key: string, unprefixed = false) {
    const pref = unprefixed ? key : this.getKeyNamePrefixed(key);
    this.store.removeItem(pref);
  }

  /**
   * Get number value
   * @param {string} key Key
   * @param {number} otherwise Default value
   * @returns {any}
   */
  getNumber(key: string, otherwise = 0): any {
    const data = this.getItem(key, otherwise);
    const num = Number(data);

    return isNaN(num) ? num : otherwise;
  }

}
