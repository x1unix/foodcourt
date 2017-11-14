import {Component, OnInit, ViewEncapsulation, ViewChild} from '@angular/core';
import {DishesService} from '../services/dishes.service';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {IDish} from '../../shared/interfaces/dish';
import {WebHelperService} from '../../shared/services/web-helper.service';
import {ResourceStatus} from '../../shared/helpers/resource-status';

const SUCCESS_ALERT_TIMEOUT = 5 * 1000; // 5 sec

@Component({
  selector: 'app-items-catalog',
  templateUrl: './items-catalog.component.html',
  styleUrls: ['./items-catalog.component.scss']
})
export class ItemsCatalogComponent extends LoadStatusComponent implements OnInit {

  /**
   * Items
   * @type {any}
   */
  items: IDish[] = null;

  /**
   * Selected items ids
   * @type {Array}
   */
  selectedIds: number[] = [];

  /**
   * Current selected item in editor
   * @type {any}
   */
  editableDish: IDish = null;

  /**
   * Show / hide editor
   * @type {boolean}
   */
  showEditor = false;

  /**
   * Disallow to close modal
   * @type {boolean}
   */
  blockModal = false;

  /**
   * Show success alert message
   * @type {boolean}
   */
  showSuccessMessage = false;

  /**
   * Show success delete message
   * @type {boolean}
   */
  showSuccessDeleteMsg = false;

  /**
   * Current editor operation
   * @type {number}
   */
  operation = 0; // 0 - Create, 1 - Update

  deleteStatus = new ResourceStatus();

  constructor(private dishes: DishesService, private web: WebHelperService) {
    super();
  }

  ngOnInit() {
    this.fetchData();
    this.deleteStatus.isIdle = true;
  }

  /**
   * Gets category name
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategory(catId: number) {
    return this.dishes.getDishCategory(catId);
  }

  /**
   * Gets category color
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategoryColor(catId: number) {
    return this.dishes.getDishCategoryColor(catId);
  }

  /**
   * Fetch items from the backend
   */
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

  /**
   * Item select event
   * @param {number} itemId
   * @param {boolean} isSelected
   */
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

  /**
   * Editor progress event
   */
  onEditLoading() {
    this.blockModal = true;
  }

  /**
   * Editor work finish event
   * @param {boolean} isSuccess Is successful
   */
  onEditFinish(isSuccess: boolean) {
    this.blockModal = false;
    this.showEditor = false;
    this.editableDish = null;

    if (isSuccess === true) {
      this.showSuccessMessage = true;

      // Hide success message after timeout finish
      setTimeout(() => this.showSuccessMessage = false, SUCCESS_ALERT_TIMEOUT);

      // Refresh data
      this.fetchData();
    }
  }

  /**
   * Create new item in the editor
   */
  openDishCreator() {
    this.editableDish = null;
    this.showEditor = true;
    this.operation = 0;
  }

  /**
   * Open item editor with the specified dish
   * @param {IDish} dish Selected dish
   */
  editItem(dish: IDish) {
    // Clear previous state
    this.editableDish = null;
    this.operation = 1;
    this.editableDish = dish;
    this.showEditor = true;
  }

  /**
   * Editor button close click event
   */
  onEditDismiss() {
    this.editableDish = null;
    this.blockModal = false;
    this.showEditor = false;
  }

  doDeleteMultiple() {
    this.showSuccessDeleteMsg = false;
    this.deleteStatus.isLoading = true;
    this.dishes.deleteMultiple(this.selectedIds).subscribe(
      () => {
        this.showSuccessDeleteMsg = true;
        this.deleteStatus.isLoaded = true;
        this.selectedIds.length = 0;

        // Hide success message after timeout finish
        setTimeout(() => this.showSuccessDeleteMsg = false, SUCCESS_ALERT_TIMEOUT);

        this.fetchData();
      }, (err) => {
        this.showSuccessDeleteMsg = false;
        this.deleteStatus.isFailed = true;
        this.deleteStatus.error = err;
      }
    );
  }

}
