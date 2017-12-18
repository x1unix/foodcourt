import { Component, OnInit } from '@angular/core';
import * as moment from 'moment';
import {UsersService} from '../shared/services';
import {ReportsService} from './reports.service';

const DAY_FRIDAY = 5;

const DISPLAYED_DATE_FORMAT = 'dddd, MMMM DD YYYY';
const SERVED_DATE_FORMAT = 'YYYYMMDD';

@Component({
  selector: 'app-reports',
  templateUrl: './reports.component.html',
  styleUrls: ['./reports.component.scss']
})
export class ReportsComponent implements OnInit {

  dateFrom: number;

  dateTill: number;

  constructor(private users: UsersService, private reports: ReportsService) { }

  ngOnInit() {
    const today = moment();
    const dayOfWeek = today.isoWeekday();

    if (dayOfWeek >= DAY_FRIDAY) {
      // Skip to next monday if today is friday or weekend
      const daysTillNextMonday = (7 - dayOfWeek) + 1;

      // Set start date from monday till tuesday
      const startDate = today.add(daysTillNextMonday, 'days');
      const endDate = startDate.add(DAY_FRIDAY - 1, 'days');

      this.dateFrom = startDate.unix();
      this.dateTill = endDate.unix();
    } else {

    }
  }

}
