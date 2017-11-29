import {Component, OnDestroy, OnInit} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import { isNil, groupBy } from 'lodash';
import * as moment from 'moment';

import { WebHelperService, MenuService, OrdersService, SessionsService } from '../../shared/services';
import { DatepickerOptions } from '../../shared/components/datepicker';
import {DishType, IDish} from '../../shared/interfaces/dish';
import {ResourceStatus} from '../../shared/helpers/resource-status';
import {LoadStatusComponent} from '../../shared/helpers/load-status-component';
import {Subscription} from 'rxjs/Rx';
import {forkJoin} from 'rxjs/observable/forkJoin';

const DISPLAYED_DATE_FORMAT = 'dddd, MMMM DD YYYY';
const SERVED_DATE_FORMAT = 'YYYYMMDD';

// Route params names
const PARAM_UID = 'userId';
const PARAM_DATE = 'date';

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

      // Init component state
      this.pickedDate = this.date.toDate();
      this.initDatePickerOptions();
      this.updateMenuAndOrder();
      this.ready = true;
    });
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

  updateMenuAndOrder() {
    this.deleteStatus.isIdle = true;
    this.saveStatus.isIdle = true;
    this.menuEmpty = true;
    this.isLoading = true;
    this.selectedIds = [];
    this.collectionChanged = false;
    this.initialSize = 0;
    this.selectedClassItems = [];

    // Fetch menu items & order data
    const menu = this.menu.getDishes(this.servedDate);
    const order = this.orders.getOrderedDishIds(this.servedDate, this.userId);

    forkJoin([menu, order]).subscribe(
      results => this.onDataLoaded.call(this, results[0], results[1]),
      error => this.onInitFail(error)
    );
  }

  onDataLoaded(menuItems: IDish[] = null, orderedIds: number[] = []) {
    this.isLoaded = true;
    this.menuEmpty = isNil(menuItems) || menuItems.length === 0;

    // Break if menu is empty
    if (this.menuEmpty) {
      return;
    }

    this.collectionChanged = false;
    this.selectedIds = isNil(orderedIds) ? [] : orderedIds;

    // Create new empty collection with empty sub-arrays for each category
    this.menuItems = this.dishTypes.map((i) => []);
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

      this.initialSize = this.selectedIds.length;

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

  }

}
