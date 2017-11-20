import { Component, OnInit } from '@angular/core';
import {DishesService} from '../services/dishes.service';
import {IDish} from '../../shared/interfaces/dish';
import {ResourceStatus} from '../../shared/helpers/resource-status';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {WebHelperService} from '../../shared/services/web-helper.service';

const ITEMS_QUERY = {
  orderBy: 'label',
  orderDir: 'desc'
};

@Component({
  selector: 'app-menu-editor',
  templateUrl: './menu-editor.component.html',
  styleUrls: ['./menu-editor.component.scss']
})
export class MenuEditorComponent extends LoadStatusComponent implements OnInit {

  dishes: IDish[] = [];
  dishesStatus = new ResourceStatus();
  saveStatus = new ResourceStatus();

  constructor(private dishesCatalogue: DishesService, private helper: WebHelperService) {
    super();
  }

  ngOnInit() {
    this.getAllDishes();
  }

  getAllDishes() {
    this.dishesStatus.isLoading = true;
    this.dishesCatalogue.getAll(ITEMS_QUERY).subscribe(
      (items: IDish[]) => {
        this.dishesStatus.isLoaded = true;
        this.dishes = items;
      }, (err) => {
        this.dishesStatus.error = this.helper.extractResponseError(err);
        this.dishesStatus.isFailed = true;
      }
    );
  }

}
