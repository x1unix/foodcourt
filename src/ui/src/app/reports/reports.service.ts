import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import { HttpParams } from '@angular/common/http';

const PARAM_FROM = 'from';
const PARAM_TILL = 'till';

@Injectable()
export class ReportsService {

  constructor(private http: HttpClient) { }

  /**
   * Gets order report for specified period
   * @param {string} from Date from
   * @param {string} till Date till
   * @returns {Observable<Object>}
   */
  getOrderStats(from: string, till: string) {
    const params = new HttpParams()
      .set(PARAM_FROM, from)
      .set(PARAM_TILL, till);

    return this.http.get('/api/orders/report', {
      params: params
    });
  }

}
