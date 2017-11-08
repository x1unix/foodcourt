import { Injectable } from '@angular/core';
import { isNil, isNumber } from 'lodash';
import { Subject } from 'rxjs/Subject';
import * as moment from 'moment';

import {LocalStorageService} from './local-storage.service';
import {IAuthSession} from '../interfaces/auth-session';
import {LoggerService} from './logger.service';

const TAG = 'SessionsService';
export const KEY_SESSION = 'session';

@Injectable()
export class SessionsService {

  private _ttl: number = null;
  private _token: string = null;
  private _authorized = false;
  private _userId: number = null;


  sessionExpired = new Subject();

  /**
   * API auth token
   * @returns {string}
   */
  get token() {
    return this._token;
  }

  /**
   * Current user ID
   * @returns {number}
   */
  get userId() {
    return this._userId;
  }

  /**
   * Is API token available (based on session)
   * @returns {boolean}
   */
  get isAuthorized() {
    return this._authorized;
  }

  // Is token alive
  get isSessionAlive() {
    if (!this._authorized) {
      return false;
    }

    if (!isNumber(this._ttl)) {
      return false;
    }

    // Session TTL in UNIX epoch time
    const ttl = this._ttl;

    // Check if ttl is not passed
    const ttlNotPassed = this.isTtlActive(ttl);

    return ttlNotPassed;
  }

  /**
   * Check if ttl period is active
   * @param {number} ttl TTL period
   * @returns {boolean}
   */
  private isTtlActive(ttl: number) {
    // Current UNIX epoch in UTC
    const now = moment.utc().unix();

    return now < ttl;
  }

  constructor(private storage: LocalStorageService, private log: LoggerService) {
    // Try to load session from the local storage
    const cached = this.tryLoadCachedSession();

    // Print debug message
    const msg = cached ? 'Session was loaded from the cache' : 'No available sessions at cache';

    this.log.info(TAG, msg);
  }

  /**
   * Try to load session from the local storage (if available and valid)
   */
  private tryLoadCachedSession() {
    // Try to get session from the local storage
    const keyExists = this.storage.hasKey(KEY_SESSION);

    // Skip if no cached session
    if (!keyExists) {
      return false;
    }

    // Try to get session object
    const session = this.storage.getObject(KEY_SESSION);

    // Skip if data is null or undefined
    if (isNil(session)) {
      return false;
    }

    // Check if Ttl is not passed
    const isTtlAlive = this.isTtlActive(session.ttl);

    if (isTtlAlive) {
      // Load session gracefully
      this.loadSession(session);
      return true;
    }

    // Purge session from storage if it's invalid
    this.purgeSession(false);
    return false;
  }

  /**
   * Load specified session
   * @param {IAuthSession} newSession Session
   */
  loadSession(newSession: IAuthSession) {
    // Assign values
    this._userId = newSession.userId;
    this._ttl = newSession.ttl;
    this._token = newSession.token;
    this._authorized = newSession.authorized;

    // Save session
    this.storage.setItem(KEY_SESSION, newSession);
  }

  /**
   * Delete session data
   * @param {boolean} emitExpiredEvent Emit SessionExpired event
   */
  purgeSession(emitExpiredEvent = false) {
    // Clear state
    this._ttl = null;
    this._userId = null;
    this._token = null;
    this._authorized = false;

    // Remove cached session
    this.storage.removeItem(KEY_SESSION);

    // Emit event (optional)
    if (emitExpiredEvent) {
      this.log.warn(TAG, 'Session expired');
      this.sessionExpired.next();
    }
  }

}
