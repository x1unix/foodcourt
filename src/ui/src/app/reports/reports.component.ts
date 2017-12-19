import { Component, OnInit } from '@angular/core';
import * as moment from 'moment';
import { isNil } from 'lodash';
import {forkJoin} from 'rxjs/observable/forkJoin';
import {UsersService, WebHelperService} from '../shared/services';
import {ReportsService} from './reports.service';
import { DateTime } from './model/date-time';
import {LoadStatusComponent} from '../shared/helpers';
import {IUser} from '../shared/interfaces/user';
import {IKeyValuePair} from '../shared/interfaces/key-value-pair';

const DAY_FRIDAY = 5;
const DAY_MONDAY = 1;

const SELECTOR_DATE_FORMAT = 'DD/MM/YYYY';
const GRID_DATE_FORMAT = 'dddd, DD MMM';
const SERVED_DATE_FORMAT = 'YYYYMMDD';

@Component({
  selector: 'app-reports',
  templateUrl: './reports.component.html',
  styleUrls: ['./reports.component.scss']
})
export class ReportsComponent extends LoadStatusComponent implements OnInit {

  dateFrom: DateTime;
  dateTill: DateTime;
  usersList: IUser[];

  reportData: any = null;
  dateLabels: IKeyValuePair<number, string>[] = [];

  initialized = false;


  constructor(private users: UsersService, private reports: ReportsService, private helper: WebHelperService) {
    super();
  }

  ngOnInit() {
    const today = moment();
    const dayOfWeek = today.isoWeekday();

    // Initialize days
    if (dayOfWeek >= DAY_FRIDAY) {
      // Skip to next monday if today is friday or weekend
      const daysTillNextMonday = (7 - dayOfWeek) + 1;

      // Set start date from monday till tuesday
      const startDate = today.add(daysTillNextMonday, 'days');
      const endDate = startDate.add(DAY_FRIDAY - 1, 'days');

      this.dateFrom = new DateTime(startDate);
      this.dateTill = new DateTime(endDate);
    } else {
      const daysTillMonday = dayOfWeek - 1;

      const startDate = today.subtract(daysTillMonday, 'days');
      const endDate = startDate.clone().add(DAY_FRIDAY - 1, 'days');

      this.dateFrom = new DateTime(startDate);
      this.dateTill = new DateTime(endDate);
    }

    this.fetchInitData();
  }

  /**
   * Generates a period list between two dates.
   * Used for report table.
   */
  buildDateLabels() {
    const start = this.dateFrom.origin.clone();
    const daysRange = start.diff(this.dateTill.origin, 'days');
    this.dateLabels = [];

    for (let d = 0; d <= daysRange; d++) {
      start.add(d, 'days');
      this.dateLabels.push({
        key: +start.format(SERVED_DATE_FORMAT),
        value: start.format(GRID_DATE_FORMAT)
      });
    }
  }

  /**
   * Fetches users list and report.
   * Used at component start.
   */
  fetchInitData() {
    this.isLoading = true;
    const users = this.users.getAll();
    const stats = this.reports.getOrderStats(this.dateFrom.served, this.dateTill.served);

    forkJoin([users, stats]).subscribe((results) => {
      this.usersList = <IUser[]> results[0];
      this.reportData = results[1];
      this.buildDateLabels();

      this.initialized = true;
      this.isLoaded = true;
    }, (err) => {
      this.error = this.helper.extractResponseError(err);
      this.isFailed = true;
    });
  }

  /**
   * Checks if user made order for specific date
   * @param {number} userId User ID
   * @param {number} date Date (YYYYMMDD)
   * @returns {boolean}
   */
  isUserOrderMade(userId: number, date: number): boolean {
    if (isNil(this.reportData[`${userId}`])) {
      return false;
    }

    return this.reportData[`${userId}`].includes(date);
  }

  fetchReport() {

  }

  retry(event: Event) {
    event.preventDefault();
    if (this.initialized) {
      this.fetchInitData();
    } else {
      this.fetchReport();
    }
    return false;
  }

}
