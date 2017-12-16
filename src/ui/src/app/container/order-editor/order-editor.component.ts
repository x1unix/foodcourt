import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import { isNil, groupBy, isObject, flatten } from 'lodash';
import * as moment from 'moment';

import { WebHelperService, MenuService, OrdersService, SessionsService } from '../../shared/services';
import { DatepickerOptions } from '../../shared/components/datepicker';
import {DishType, IDish} from '../../shared/interfaces/dish';
import {ResourceStatus} from '../../shared/helpers/resource-status';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {Subscription} from 'rxjs/Rx';
import {forkJoin} from 'rxjs/observable/forkJoin';
import {DatepickerComponent} from '../../shared/components/datepicker/datepicker.component';
import {ILockStatus} from '../../shared/interfaces/lock-status';

const DISPLAYED_DATE_FORMAT = 'dddd, MMMM DD YYYY';
const SERVED_DATE_FORMAT = 'YYYYMMDD';

// Route params names
const PARAM_UID = 'userId';
const PARAM_DATE = 'date';

const MSG_CONFIRM = 'Do you want to move to other date without saving changes? All unsaved changes will be lost.';

@Component({
  selector: 'app-order-editor',
  templateUrl: './order-editor.component.html',
  styleUrls: ['./order-editor.component.scss']
})
export class OrderEditorComponent extends LoadStatusComponent implements OnInit, OnDestroy {

  /**
   * Is component ready to work
   * @type {boolean}
   */
  ready = false;

  /**
   * User menu id
   * @type {string}
   */
  userId: string = null;

  /**
   * Array of ordered items
   * @type {Array}
   */
  selectedIds: number[] = [];

  /**
   * Grouped list of dishes in the current menu
   * @type {Array}
   */
  menuItems: IDish[][] = [];

  /**
   * Is menu empty
   * @type {boolean}
   */
  menuEmpty = true;

  /**
   * Is order editable
   * @type {boolean}
   */
  orderEditable = true;

  /**
   * Menu save progress status
   * @type {ResourceStatus}
   */
  saveStatus = new ResourceStatus();

  /**
   * Menu delete progress status
   * @type {ResourceStatus}
   */
  deleteStatus = new ResourceStatus();

  /**
   * Route params change subscription
   * @type {Subscription}
   */
  routeParams$: Subscription = null;

  /**
   * Date to be displayed on UI
   * @type {string}
   */
  displayedDate: string = null;

  /**
   * Initial ordered items count
   * @type {number}
   */
  initialSize = 0;

  /**
   * Is collection changed
   * @type {boolean}
   */
  collectionChanged = false;

  /**
   * ng2-datepicker options
   * @type {any}
   */
  datePickerOptions: DatepickerOptions = null;

  pickedDate: Date = null;

  selectedClassItems: number[] = [];

  orderStructError = '';

  asDifferentUser = false;

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
    return this.menu.dishTypes;
  }

  get orderEmpty(): boolean {
    return this.selectedIds.length === 0;
  }

  constructor(
    private helper: WebHelperService,
    private menu: MenuService,
    private orders: OrdersService,
    private session: SessionsService,
    private route: ActivatedRoute
  ) {
    super();
  }

  ngOnInit() {
    this.routeParams$ = this.route.params.subscribe(params => {
      // Read route params (date & user id) (optional)
      this.userId = isNil(params[PARAM_UID]) ? this.session.userId : params[PARAM_UID];
      this.date = isNil(params[PARAM_DATE]) ? moment().utc() : moment(params[PARAM_DATE], SERVED_DATE_FORMAT);

      // Check user mode
      this.asDifferentUser = `${this.userId}` !== `${this.session.userId}`;

      // Init component state
      this.pickedDate = this.date.toDate();
      this.initDatePickerOptions();
      this.updateMenuAndOrder();
      this.ready = true;
    });
  }

  private detectOrderEditable() {

  }

  private setSelectedCategoryValue(catId, value: number = null) {
    this.selectedClassItems[catId] = value;
  }

  /**
   * Dish select event handler
   * @param {boolean} isChecked is element checked
   * @param {number} category dish category id
   * @param {number} itemId dish id
   */
  onItemCheckToggle(isChecked: boolean, category: number, itemId: number) {
    this.collectionChanged = true;

    if (!isChecked) {
      // Just remove item if it's unchecked
      this.setSelectedCategoryValue(category, null);
      return;
    }

    // Remove main+garnish if special is selected and vice versa
    switch (category) {
      case DishType.special:
        this.setSelectedCategoryValue(DishType.garnish, null);
        this.setSelectedCategoryValue(DishType.main, null);
        break;
      case DishType.main:
        this.setSelectedCategoryValue(DishType.special, null);
        break;
      case DishType.garnish:
        this.setSelectedCategoryValue(DishType.special, null);
        break;
      default:
        break;
    }

    // Update selected category
    this.setSelectedCategoryValue(category, itemId);
  }

  ngOnDestroy() {
    this.routeParams$.unsubscribe();
    this.menuItems = null;
    this.selectedIds = null;
    this.selectedDate = null;
    this.selectedClassItems = null;
  }

  initDatePickerOptions() {
    const currentYear = this.date.year();
    this.datePickerOptions = {
      minYear: currentYear,
      firstCalendarDay: 1,
      displayFormat: SERVED_DATE_FORMAT
    };
  }

  openPicker(picker: DatepickerComponent) {
    if ((!picker.isOpened) && (this.collectionChanged === true)) {
      const confirm = window.confirm(MSG_CONFIRM);

      if (!confirm) {
        return;
      }
    }

    picker.toggle();
  }

  updateMenuAndOrder() {
    this.deleteStatus.isIdle = true;
    this.saveStatus.isIdle = true;
    this.menuEmpty = true;
    this.isLoading = true;
    this.selectedIds = [];
    this.collectionChanged = false;
    this.orderEditable = true;
    this.initialSize = 0;
    this.selectedClassItems = [];
    this.orderStructError = '';

    // Fetch menu items & order data
    const menu = this.menu.getDishes(this.servedDate);
    const order = this.orders.getOrderedDishIds(this.servedDate, this.userId);
    const status = this.menu.getMenuStatus(this.servedDate);

    forkJoin([menu, order, status]).subscribe(
      results => this.onDataLoaded.apply(this, results),
      error => this.onInitFail(error)
    );
  }

  onDataLoaded(menuItems: IDish[] = null, orderedIds: number[] = [], menuStatus: ILockStatus) {
    this.isLoaded = true;
    this.menuEmpty = isNil(menuItems) || menuItems.length === 0;

    this.orderEditable = isObject(menuStatus) && (menuStatus.locked === false);

    // Break if menu is empty
    if (this.menuEmpty) {
      return;
    }

    this.collectionChanged = false;
    this.selectedIds = isNil(orderedIds) ? [] : orderedIds;

    // Create new empty collection with empty sub-arrays for each category
    this.menuItems = this.dishTypes.map((i) => []);
    this.selectedClassItems.length = 0;
    this.selectedClassItems.length = this.dishTypes.length;

    if (isNil(orderedIds)) {
      orderedIds = [];
    }

    // Group by class and fill the collection
    if (!isNil(menuItems)) {
      // Fill selected ids list and categories
      this.selectedIds = menuItems.map((i) => {

        // Put in selected cat items
        if (orderedIds.includes(i.id)) {
          this.selectedClassItems[i.type] = i.id;
        }

        return i.id;
      });

      this.initialSize = orderedIds.length;

      const grouped = groupBy(menuItems, 'type');
      Object.keys(grouped).forEach((groupId) => {
        this.menuItems[groupId] = [...grouped[groupId]];
      });
    }
  }

  /**
   * Menu or order list load error handler
   * @param error
   */
  onInitFail(error) {
    this.isFailed = true;
    this.error = this.helper.extractResponseError(error);
  }

  /**
   * Gets category name
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategory(catId: number) {
    return this.menu.getDishCategory(catId);
  }

  /**
   * Gets category color
   * @param {number} catId Category id
   * @returns {string}
   */
  getDishCategoryColor(catId: number) {
    return this.menu.getDishCategoryColor(catId);
  }

  /**
   * Picked date change event handler
   * @param {Date} newDate
   */
  onDateChange(newDate: Date) {
    this.date = moment(newDate);
    this.updateMenuAndOrder();
  }

  saveChanges() {
    this.orderStructError = '';

    // Just delete items if selected collection is empty
    if (this.selectedClassItems.length === 0) {
      return this.deleteItems();
    }

    let structOk, structErr;

    [structOk, structErr] = this.checkMenuStructure();

    if (!structOk) {
      this.orderStructError = structErr;
      return;
    }

    const dishes = flatten(this.selectedClassItems);

    this.saveStatus.isLoading = true;
    this.orders.orderDishes(dishes, this.servedDate, this.userId).subscribe(
      () => this.onSaveSuccess(),
      (err) => this.onSaveFail(err)
    );
  }

  deleteItems() {
    this.deleteStatus.isLoading = true;
    this.orders.deleteOrder(this.servedDate, this.userId).subscribe(
      () => this.onDeleteSuccess(),
      (err) => this.onDeleteFail(err)
    );
  }

  onDeleteSuccess() {
    this.collectionChanged = false;

    // Rebuild selected sorted items array
    this.selectedClassItems.length = 0;
    this.selectedClassItems.length = this.dishTypes.length;

    this.initialSize = 0;
    this.deleteStatus.isLoaded = true;
    setTimeout(() => this.deleteStatus.isIdle = true, 3000);
  }

  onDeleteFail(error) {
    this.deleteStatus.error = this.helper.extractResponseError(error);
    this.deleteStatus.isFailed = true;
  }

  onSaveSuccess() {
    this.collectionChanged = false;
    this.initialSize = flatten(this.selectedClassItems).length;
    this.saveStatus.isLoaded = true;
    setTimeout(() => this.saveStatus.isIdle = true, 3000);
  }

  onSaveFail(error) {
    this.saveStatus.error = this.helper.extractResponseError(error);
    this.saveStatus.isFailed = true;
  }

  /**
   * Check if the main dish + garnish pair was selected correctly
   * @returns {[boolean , string]}
   */
  private checkMenuStructure(): [boolean, string] {
    const selected = this.selectedClassItems;
    if (!isNil(selected[DishType.garnish]) && isNil(selected[DishType.main])) {
      return [false, 'Please select the main dish'];
    }

    if (!isNil(selected[DishType.main]) && isNil(selected[DishType.garnish])) {
      return [false, 'Please select the garnish'];
    }

    return [true, ''];
  }

}
