import { Component, OnInit } from '@angular/core';
import { isNil, groupBy } from 'lodash';
import * as moment from 'moment';

import {DishesService} from '../services/dishes.service';
import {IDish} from '../../shared/interfaces/dish';
import {ResourceStatus} from '../../shared/helpers/resource-status';
import {WebHelperService, MenuService} from '../../shared/services';
import { DropEvent } from '../../shared/interfaces/drop-event';

const ITEMS_QUERY = {
  orderBy: 'label',
  orderDir: 'asc'
};

const DISPLAYED_DATE_FORMAT = 'dddd, MMMM DD YYYY';
const SERVED_DATE_FORMAT = 'YYYYMMDD';

/**
 * Menu editor page component
 */
@Component({
  selector: 'app-menu-editor',
  templateUrl: './menu-editor.component.html',
  styleUrls: ['./menu-editor.component.scss']
})
export class MenuEditorComponent implements OnInit {

  /**
   * List of all available dishes
   * @type {Array}
   */
  dishes: IDish[] = [];

  /**
   * List of dishes in menu
   * @type {Array}
   */
  menuItems: IDish[][] = [];

  /**
   * Dishes list fetch progress status
   * @type {ResourceStatus}
   */
  dishesStatus = new ResourceStatus();

  /**
   * Menu save progress status
   * @type {ResourceStatus}
   */
  saveStatus = new ResourceStatus();

  /**
   * Menu items fetch status
   * @type {ResourceStatus}
   */
  menuItemsStatus = new ResourceStatus();

  /**
   * Date to be displayed on UI
   * @type {any}
   */
  displayedDate: string = null;

  /**
   * List of selected id's
   * @type {Array}
   */
  selectedIds: number[] = [];

  catalogDropZone = ['allItemsZone'];

  cartDropZone = ['selectedItemsZone'];

  dragTrashStart = false;

  /**
   * Date to be send on server
   * @type {any}
   */
  private servedDate: string = null;

  /**
   * Selected date in form (private)
   * @type {any}
   */
  private selectedDate: moment.Moment = null;

  /**
   * Current selected date
   * @returns {moment.Moment}
   */
  get date(): moment.Moment {
    return this.selectedDate;
  }

  set date(newDate: moment.Moment) {
    this.selectedDate = newDate;
    this.displayedDate = newDate.format(DISPLAYED_DATE_FORMAT);
    this.servedDate = newDate.format(SERVED_DATE_FORMAT);
  }

  get dishTypes() {
    return this.dishesCatalogue.dishTypes;
  }

  get menuEmpty(): boolean {
    return this.selectedIds.length === 0;
  }

  constructor(private dishesCatalogue: DishesService, private helper: WebHelperService, private menu: MenuService) {}

  ngOnInit() {
    this.date = moment().utc();
    this.getAllDishes();
    this.updateMenuItemsList();
  }

  /**
   * Gets category name
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategory(catId: number) {
    return this.dishesCatalogue.getDishCategory(catId);
  }

  /**
   * Gets category color
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategoryColor(catId: number) {
    return this.dishesCatalogue.getDishCategoryColor(catId);
  }

  updateMenuItemsList() {
    this.menuItemsStatus.isLoading = true;
    this.selectedIds = [];
    this.menu.getDishes(this.servedDate)
      .subscribe((i) => this.onMenuItemsFetch(i), (e) => this.onMenuItemsFail(e));
  }

  onMenuItemsFetch(items: IDish[] = null) {
    // Set load status
    this.menuItemsStatus.isLoaded = true;

    // Create new empty collection with empty sub-arrays for each category
    this.menuItems = this.dishTypes.map((i) => []);

    // Group by class and fill the collection
    if (!isNil(items)) {
      // Fill selected ids list
      this.selectedIds = items.map((i) => i.id);

      const grouped = groupBy(items, 'type');
      Object.keys(grouped).forEach((groupId) => {
        this.menuItems[groupId] = [...grouped[groupId]];
      });
    }
  }

  onMenuItemsFail(error) {
    this.menuItemsStatus.isFailed = true;
    this.menuItemsStatus.error = this.helper.extractResponseError(error);
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

  onDrop(data: DropEvent<IDish>) {
    const dish = data.dragData;
    this.menuItems[dish.type].push(dish);
    this.selectedIds.push(dish.id);
  }

  onItemRemove(dish: IDish) {
    this.dragTrashStart = false;
    const src = this.menuItems[dish.type];

    if (isNil(src)) {
      return;
    }

    const index = src.indexOf(dish);

    if (index === -1) {
      return;
    }

    src.splice(index, 1);

    // remove from selected id's
    this.selectedIds.splice(this.selectedIds.indexOf(dish.id), 1);
  }

}
