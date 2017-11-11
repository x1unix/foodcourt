import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {IDish, DISH_TYPES, DISH_TYPE_COLORS} from '../../shared/interfaces/dish';
import {Observable} from 'rxjs/Observable';

@Injectable()
export class DishesService {

  constructor(private http: HttpClient) { }

  getAll(): Observable<IDish[]> {
    return <Observable<IDish[]>> this.http.get('/api/dishes');
  }

  getDishCategory(categoryId: number): string {
    return DISH_TYPES[categoryId];
  }

  getDishCategoryColor(catId: number): string {
    return DISH_TYPE_COLORS[catId];
  }

}
