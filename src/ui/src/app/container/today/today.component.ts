import {Component, OnDestroy, OnInit} from '@angular/core';
import { isNil } from 'lodash';
import * as moment from 'moment';
import {LoadStatusComponent} from '../../shared/helpers';
import {MenuService, OrdersService, SessionsService, WebHelperService} from '../../shared/services';
import {IDish} from '../../shared/interfaces/dish';

const DATE_FORMAT = 'YYYYMMDD';
const DISPLAYED_DATE_FORMAT = 'dddd, MMMM DD YYYY';

@Component({
  selector: 'app-today',
  templateUrl: './today.component.html',
  styleUrls: ['./today.component.scss']
})
export class TodayComponent extends LoadStatusComponent implements OnInit, OnDestroy {

  /**
   * Is order list is empty
   * @type {boolean}
   */
  orderEmpty = false;

  /**
   * List of dishes
   * @type {IDish[]}
   */
  dishes: IDish[] = [];

  /**
   * Displayed date
   * @type {null}
   */
  displayedDate: string = null;

  /**
   * Current user id
   * @type {null}
   */
  private userId: number = null;

  /**
   * Order date
   * @type {null}
   */
  private orderDate: string = null;

  constructor(
    private orders: OrdersService,
    private menu: MenuService,
    private session: SessionsService,
    private helper: WebHelperService
  ) {
    super();
  }

  ngOnInit() {
    this.userId = this.session.userId;
    this.orderDate = moment().format(DATE_FORMAT);
    this.displayedDate = moment().format(DISPLAYED_DATE_FORMAT);

    this.fetchOrder();
  }

  ngOnDestroy() {
    this.userId = null;
    this.orderDate = null;
    this.orders = null;
    this.session = null;
    this.helper = null;
    this.dishes = null;
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
   * Fetch order info
   */
  fetchOrder() {
    this.isLoading = true;
    this.dishes = [];

    this.orders.getOrderedDishes(this.orderDate, this.userId).subscribe(
      (data) => this.onFetchSuccess(data),
      (err) => this.onFetchFail(err)
    );

  }

  /**
   * Fetch success event handler
   * @param {IDish[]} dishes
   */
  onFetchSuccess(dishes: IDish[] = null) {
    this.isLoaded = true;

    this.orderEmpty = isNil(dishes) || (dishes.length === 0);

    if (!this.orderEmpty) {
      this.dishes = dishes;
    }
  }

  /**
   * Fetch fail event handler
   * @param err
   */
  onFetchFail(err) {
    this.isFailed = true;
    this.error = this.helper.extractResponseError(err);
  }

}
