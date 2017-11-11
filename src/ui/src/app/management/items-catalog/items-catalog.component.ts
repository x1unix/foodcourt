import {Component, OnInit, ViewEncapsulation} from '@angular/core';
import {DishesService} from '../services/dishes.service';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {IDish} from '../../shared/interfaces/dish';
import {WebHelperService} from '../../shared/services/web-helper.service';

@Component({
  selector: 'app-items-catalog',
  templateUrl: './items-catalog.component.html',
  styleUrls: ['./items-catalog.component.scss'],
  providers: [
    DishesService
  ]
})
export class ItemsCatalogComponent extends LoadStatusComponent implements OnInit {

  items: IDish[] = null;

  constructor(private dishes: DishesService, private web: WebHelperService) {
    super();
  }

  ngOnInit() {
    this.fetchData();
  }

  getDishCategory(catId: number) {
    return this.dishes.getDishCategory(catId);
  }

  getDishCategoryColor(catId: number) {
    return this.dishes.getDishCategoryColor(catId);
  }

  fetchData() {
    this.isLoading = true;
    this.dishes.getAll().subscribe(
      (data) => {
        this.items = data;
        this.isLoaded = true;
      }, (err) => {
        this.error = this.web.extractResponseError(err);
        this.isFailed = true;
      }
    );
  }

}
