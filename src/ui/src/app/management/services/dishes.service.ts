import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {IDish, DISH_TYPES, DISH_TYPE_COLORS} from '../../shared/interfaces/dish';
import {Observable} from 'rxjs/Observable';
import {IMessage} from '../../shared/interfaces/message';
import {RequestOptions} from '@angular/http';

@Injectable()
export class DishesService {

  constructor(private http: HttpClient) { }

  /**
   * Get list of all dishes
   * @returns {Observable<IDish[]>}
   */
  getAll(): Observable<IDish[]> {
    return <Observable<IDish[]>> this.http.get('/api/dishes');
  }

  /**
   * Get dish category name
   * @param {number} categoryId Category index
   * @returns {string}
   */
  getDishCategory(categoryId: number): string {
    return DISH_TYPES[categoryId];
  }

  /**
   * Get dish category color
   * @param {number} catId Category index
   * @returns {string}
   */
  getDishCategoryColor(catId: number): string {
    return DISH_TYPE_COLORS[catId];
  }

  /**
   * Create a new dish
   * @param {IDish} dish Dish
   * @returns {Observable<IMessage>}
   */
  addDish(dish: IDish): Observable<IMessage> {
    dish.type = Number(dish.type);
    return <Observable<IMessage>> this.http.post('/api/dishes', dish);
  }

  /**
   * Update an existing dish
   * @param {IDish} dish Dish
   * @returns {Observable<IMessage>}
   */
  updateDish(dish: IDish): Observable<IMessage> {
    dish.type = Number(dish.type);
    return <Observable<IMessage>> this.http.put(`/api/dishes/${dish.id}`, dish);
  }

  /**
   * Delete multiple items by id
   * @param {string[]} dishIds Items ids
   * @returns {Observable<IMessage>}
   */
  deleteMultiple(dishIds: number[]): Observable<IMessage> {
    return <Observable<IMessage>> this.http.post('/api/dishes/purge', dishIds);
  }

}
