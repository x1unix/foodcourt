import { Injectable } from '@angular/core';
import { HttpClient} from '@angular/common/http';

import {IAuthCredentials} from '../interfaces/auth-credentials';
import {Observable} from 'rxjs/Observable';
import {IMessage} from '../interfaces/message';
import {IAuthSession} from '../interfaces/auth-session';

export const URL_LOGIN = '/api/login';

@Injectable()
export class AuthService {

  constructor(private http: HttpClient) { }

  /**
   * Get session data
   * @returns {Observable<IAuthSession>}
   */
  getSessionStatus(): Observable<IAuthSession> {
    return <Observable<IAuthSession>> this.http.get('/api/session');
  }

  /**
   * Try to login
   * @param {IAuthCredentials} credentials Credentials
   * @returns {Observable<IAuthSession>}
   */
  login(credentials: IAuthCredentials): Observable<IAuthSession> {
    return <Observable<IAuthSession>> this.http.post(URL_LOGIN, credentials);
  }

  /**
   * Logout
   * @returns {Observable<IMessage>}
   */
  logout(): Observable<IMessage> {
    return <Observable<IMessage>> this.http.post('/api/logout', null);
  }

}
