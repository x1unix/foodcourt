import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {IUser} from '../interfaces/user';

const GROUP_NAMES = [
  'Administrator',
  'Manager',
  'Customer'
];

@Injectable()
export class UsersService {

  constructor(private http: HttpClient) { }

  /**
   * List of user groups
   * @returns {string[]}
   */
  get groups() {
    return GROUP_NAMES;
  }

  /**
   * Gets group name
   * @param {number} grpId Group ID
   * @returns {string}
   */
  getGroupName(grpId: number): string {
    return GROUP_NAMES[grpId];
  }

  /**
   * Gets list of all users
   * @returns {Observable<IUser[]>}
   */
  getAll(): Observable<IUser[]> {
    return <Observable<IUser[]>> this.http.get(`/api/users`);
  }

}
