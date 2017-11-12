import {Component, OnInit, ViewEncapsulation, ViewChild} from '@angular/core';
import {DishesService} from '../services/dishes.service';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {IDish} from '../../shared/interfaces/dish';
import {WebHelperService} from '../../shared/services/web-helper.service';

@Component({
  selector: 'app-items-catalog',
  templateUrl: './items-catalog.component.html',
  styleUrls: ['./items-catalog.component.scss']
})
export class ItemsCatalogComponent extends LoadStatusComponent implements OnInit {

  items: IDish[] = null;

  selectedIds: number[] = [];

  editableDish: IDish = null;

  // Views
  showEditor = false;

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

  onItemSelect(itemId: number, isSelected: boolean) {
    if (isSelected) {
      this.selectedIds.push(itemId);
    } else {
      const itemIdIndex = this.selectedIds.indexOf(itemId);

      if (itemIdIndex === -1) {
        return;
      }

      this.selectedIds.splice(itemIdIndex, 1);
    }
  }

  openDishCreator() {
    this.showEditor = true;
  }

}
