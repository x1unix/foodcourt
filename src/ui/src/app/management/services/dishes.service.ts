import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {IDish} from '../../shared/interfaces/dish';
import {Observable} from 'rxjs/Observable';

@Injectable()
export class DishesService {

  constructor(private http: HttpClient) { }

  getAll(): Observable<IDish[]> {
    return <Observable<IDish[]>> this.http.get('/api/dishes');
  }

}
