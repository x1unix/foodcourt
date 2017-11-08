import {Component, OnInit, ViewEncapsulation} from '@angular/core';
import * as moment from 'moment';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.scss'],
  encapsulation: ViewEncapsulation.None
})
export class DashboardComponent implements OnInit {

  todayString = '';

  constructor() { }

  ngOnInit() {
    this.todayString = moment().format('dddd, MMMM Do YYYY');
  }

}
