import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs/Rx';
import {IDish} from '../interfaces/dish';

@Injectable()
export class MenuService {

  constructor(private http: HttpClient) { }

  /**
   * Get list of dishes in menu for a specific date.
   *
   * @param {string} date Menu date (date format must be YYYYMMDD (for ex. 20171121))
   * @returns {Observable<IDish[]>}
   */
  getDishes(date: string): Observable<IDish[]> {
    return <Observable<IDish[]>> this.http.get(`/api/menu/${date}/dishes`);
  }

}
