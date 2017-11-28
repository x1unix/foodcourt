import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';

@Injectable()
export class OrdersService {

  constructor(private http: HttpClient) { }

  /**
   * Order dishes for specific user at specific date
   * @param {number[]} dishes Dishes (ids)
   * @param {string} date Date (YYYYMMDD)
   * @param {number} userId User ID
   * @returns {Observable<number[]>}
   */
  orderDishes(dishes: number[], date: string, userId: number): Observable<number[]> {
    return <Observable<number[]>> this.http.post(`/api/orders/${date}/users/${userId}`, dishes);
  }

}
