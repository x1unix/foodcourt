import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs/Observable';
import {IMessage} from '../interfaces/message';
import {IDish} from '../interfaces/dish';

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
   * Delete order for specific user at specific date
   * @param {string} date Date (YYYYMMDD)
   * @param {number} userId User ID
   * @returns {Observable<IMessage>}
   */
  deleteOrder(date: string, userId: string): Observable<IMessage> {
    return <Observable<IMessage>> this.http.delete(`/api/orders/${date}/users/${userId}`);
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

  /**
   * Get ordered dishes by user at specific date
   * @param {string} date Date (YYYYMMDD)
   * @param {number} userId User id
   * @returns {Observable<IDish>}
   */
  getOrderedDishes(date: string, userId: number): Observable<IDish[]> {
    return <Observable<IDish[]>> this.http.get(`/api/orders/${date}/users/${userId}/dishes`);
  }

}
