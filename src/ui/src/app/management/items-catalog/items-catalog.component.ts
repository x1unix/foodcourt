import {Component, OnInit, ViewEncapsulation, ViewChild, OnDestroy} from '@angular/core';
import {DishesService} from '../services/dishes.service';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {IDish} from '../../shared/interfaces/dish';
import {WebHelperService} from '../../shared/services/web-helper.service';
import {ResourceStatus} from '../../shared/helpers/resource-status';
import {FormControl, FormGroup} from '@angular/forms';
import {SearchQuery} from '../helpers/search-query';
import {Subscription} from 'rxjs/Rx';

const SUCCESS_ALERT_TIMEOUT = 5 * 1000; // 5 sec
const FORM_UPD_PARAMS = {onlySelf: false, emitEvent: false};
const FORM_SUBMIT_TIMEOUT = 500; // 1s

const ORDER_BY = [
  {
    label: 'added (ASC)',
    key: 'id',
    dir: 'asc'
  },
  {
    label: 'added (DESC)',
    key: 'id',
    dir: 'desc'
  },
  {
    label: 'name (ASC)',
    key: 'label',
    dir: 'asc'
  },
  {
    label: 'name (DESC)',
    key: 'label',
    dir: 'desc'
  },
  {
    label: 'type (DESC)',
    key: 'type',
    dir: 'desc'
  },
  {
    label: 'type (ASC)',
    key: 'type',
    dir: 'asc'
  }
];

@Component({
  selector: 'app-items-catalog',
  templateUrl: './items-catalog.component.html',
  styleUrls: ['./items-catalog.component.scss']
})
export class ItemsCatalogComponent extends LoadStatusComponent implements OnInit, OnDestroy {

  /**
   * Items
   * @type {any}
   */
  items: IDish[] = [];

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

  /**
   * Delete operation progress
   * @type {ResourceStatus}
   */
  deleteStatus = new ResourceStatus();

  /**
   * Search filter form
   * @type {any}
   */
  searchBar: FormGroup = null;

  /**
   * Search params
   * @type {any}
   */
  searchParams: SearchQuery = null;

  /**
   * Current data offset. Used to get next chunk
   * @type {number}
   */
  offset = 0;

  /**
   * Size of the last data chunk. Used to predict if more items are available
   * @type {number}
   */
  lastChunkSize = 0;

  formUpdate$: Subscription = null;

  get orderList() {
    return ORDER_BY;
  }

  itemsPerPage = 32;

  resetDone = true;

  constructor(private dishes: DishesService, private web: WebHelperService) {
    super();
  }

  ngOnInit() {

    this.deleteStatus.isIdle = true;

    this.searchBar = new FormGroup({
      order: new FormControl(0),
      query: new FormControl('')
    });

    // Subscribe for filters change with timeout
    this.formUpdate$ = this.searchBar.valueChanges.
      debounceTime(FORM_SUBMIT_TIMEOUT).
      subscribe(() => this.onFilterUpdate());

    this.updateSearchParams();
    this.fetchData();
  }

  ngOnDestroy() {
    this.formUpdate$.unsubscribe();
  }

  /**
   * Enable or disable search bar
   * @param isEnabled
   */
  setSearchBarState(isEnabled) {
    if (isEnabled) {
      this.searchBar.enable(FORM_UPD_PARAMS);
    } else {
      this.searchBar.disable(FORM_UPD_PARAMS);
    }
  }

  /**
   * Generate new search params
   */
  updateSearchParams() {
    const values = this.searchBar.value;
    const selectedOrder = ORDER_BY[values.order];

    this.searchParams = {
      orderBy: selectedOrder.key,
      orderDir: selectedOrder.dir,
      searchQuery: values.query,
      offset: this.offset,
      limit: this.itemsPerPage
    };
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
    this.setSearchBarState(false);

    this.dishes.getAll(this.searchParams).subscribe(
      (data: IDish[]) => {
        this.resetDone = false;
        this.items = [...this.items, ...data];
        this.offset = this.items.length;
        this.lastChunkSize = data.length;
        this.isLoaded = true;
        this.setSearchBarState(true);
      }, (err) => {
        this.error = this.web.extractResponseError(err);
        this.isFailed = true;
        this.setSearchBarState(true);
      }
    );
  }

  getNextChunk() {
    if (this.isLoading) {
      return;
    }

    this.updateSearchParams();
    this.fetchData();
  }

  /**
   * Filter update event handler
   */
  onFilterUpdate() {
    this.resetDone = true;
    this.offset = 0; // Clear offset position
    this.lastChunkSize = 0; // Clear last chunk size
    this.items.length = 0; // The same as "this.items = [];"
    this.updateSearchParams(); // update search params
    this.fetchData(); // Fetch data from the API
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
      this.onFilterUpdate();
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

        this.onFilterUpdate();
      }, (err) => {
        this.showSuccessDeleteMsg = false;
        this.deleteStatus.isFailed = true;
        this.deleteStatus.error = err;
      }
    );
  }

}
