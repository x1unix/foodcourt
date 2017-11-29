import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {IMessage} from '../interfaces/message';

/**
 * Orders management service
 */
@Injectable()
export class OrdersService {

  constructor(private http: HttpClient) { }

  /**
   * Order dishes for specific user at specific date
   * @param {number[]} dishes Dishes (ids)
   * @param {string} date Date (YYYYMMDD)
   * @param {number} userId User ID
   * @returns {Observable<IMessage>}
   */
  orderDishes(dishes: number[], date: string, userId: string): Observable<IMessage> {
    return <Observable<IMessage>> this.http.post(`/api/orders/${date}/users/${userId}`, dishes);
  }

  /**
   * Get list of ordered dishes
   * @param {string} date Date (YYYYMMDD)
   * @param {string} userId User ID
   * @returns {Observable<number[]>}
   */
  getOrderedDishIds(date: string, userId: string): Observable<number[]> {
    return <Observable<number[]>> this.http.get(`/api/orders/${date}/users/${userId}`);
  }

}
